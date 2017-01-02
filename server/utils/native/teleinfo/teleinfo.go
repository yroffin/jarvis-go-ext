package teleinfo

import (
	"os"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
)

// On linux
// apply: stty 1200 cs7 evenp cstopb -igncr -inlcr -brkint -icrnl -opost -isig -icanon -iexten -F /dev/ttyUSB0
// to configure teleinfo tty
// Cf. https://hallard.me/gestion-de-la-teleinfo-avec-un-raspberry-pi-et-une-carte-arduipi/

// Teleinfo : instance Teleinfo device struct
type Teleinfo struct {
	entries map[string]string
}

// Trame EDF
type TeleinfoTrame struct {
	etiquette string // ETIQUETTE (4 à 8 caractères)
	data      string // DATA (1 à 12 caractères)
	hecksum   string // CHECKSUM (caractère de contrôle : Somme de caractère)
}

var instance *Teleinfo
var once sync.Once

// GetInstance : singleton instance
func GetInstance() *Teleinfo {
	once.Do(func() {
		instance = new(Teleinfo)
		instance.init()
	})
	return instance
}

// declare canal
var canal = make(chan byte, 5)

// handleReadFile : read file
func handleReadFile(device string) error {

	logger.NewLogger().WithFields(logrus.Fields{
		"device": device,
	}).Info("handleReadFile")

	s, _ := os.OpenFile(device, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)

	// Receive reply
	for {
		buf := make([]byte, 128)
		var len, err = s.Read(buf)
		if err != nil { // err will equal io.EOF
			break
		}
		for i := 0; i < len; i++ {
			canal <- buf[i]
		}
	}

	logger.NewLogger().WithFields(logrus.Fields{
		"status": "done",
	}).Info("handleReadFile")

	return nil
}

/**
* phase1: LF (0x0A)
* phase2: ETIQUETTE (4 à 8 caractères)
* phase3: SP (0x20)
* phase4: DATA (1 à 12 caractères)
* phase5: SP (0x20)
* phase6: CHECKSUM (caractère de contrôle : Somme de caractère)
* CR (0x0D)
 */

// submit trame
func submit(teleinfo *Teleinfo, trame TeleinfoTrame) {
	teleinfo.entries[trame.etiquette] = trame.data
}

// single trame
func handleTrame(teleinfo *Teleinfo, trame string) {
	var espace int
	var send TeleinfoTrame
	for i := 0; i < len(trame); i++ {
		switch {
		case trame[i] == 0x20:
			espace++
			continue
		default:
			if espace == 0 {
				send.etiquette += string([]byte{trame[i]})
			}
			if espace == 1 {
				send.data += string([]byte{trame[i]})
			}
			continue
		}
	}
	// submit new value
	submit(teleinfo, send)
}

// all trames detection
func handleTrames(teleinfo *Teleinfo, trame string) {
	var send string
	for i := 0; i < len(trame); i++ {
		switch {
		case trame[i] == 0x0A:
			send = ""
			continue
		case trame[i] == 0x0D:
			handleTrame(teleinfo, send)
			continue
		default:
			send += string([]byte{trame[i]})
		}
	}
}

// worker to consume file
func worker(teleinfo *Teleinfo) {
	var trame string
	var etx bool

	// wait for ETX 0x003
	for i := 0; etx == false; i++ {
		var x = <-canal
		if x == 0x03 {
			etx = true
		}
	}

	// daemon
	for {
		var x = <-canal
		switch {
		case x == 0x03:
			// wait for ETX 0x003
			handleTrames(teleinfo, trame)
			break
		case x == 0x02:
			// wait for STX 0x002
			trame = ""
			break
		case x != 0x02 && x != 0x03:
			// other
			trame += string([]byte{x})
			break
		}
	}
}

var mapMutex sync.Once

// get values
func (teleinfo *Teleinfo) GetEntries(entries map[string]string) map[string]string {
	for key, value := range teleinfo.entries {
		entries[key] = value
	}
	return entries
}

// Init : Init
func (teleinfo *Teleinfo) init() {
	// add map
	teleinfo.entries = make(map[string]string)

	// start worker
	go handleReadFile(getTeleinfoFile())
	go worker(teleinfo)
	logger.NewLogger().WithFields(logrus.Fields{}).Info("Init ok")
}
