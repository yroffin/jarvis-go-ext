/**
 * Copyright 2017 Yannick Roffin
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package teleinfo

import (
	"os"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
	log "github.com/yroffin/jarvis-go-ext/logger"
)

// On linux
// apply: stty 1200 cs7 evenp cstopb -igncr -inlcr -brkint -icrnl -opost -isig -icanon -iexten -F /dev/ttyUSB0
// to configure teleinfo tty
// Cf. https://hallard.me/gestion-de-la-teleinfo-avec-un-raspberry-pi-et-une-carte-arduipi/

// Teleinfo : instance Teleinfo device struct
type Teleinfo struct {
	Entries map[string]string
}

// Trame EDF
type TeleinfoTrame struct {
	etiquette string // ETIQUETTE (4 à 8 caractères)
	data      string // DATA (1 à 12 caractères)
	hecksum   string // CHECKSUM (caractère de contrôle : Somme de caractère)
}

var instance *Teleinfo
var once sync.Once
var mutex = &sync.Mutex{}

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

	logrus.WithFields(logrus.Fields{
		"device": device,
	}).Info("handleReadFile")

	s, err := os.OpenFile(device, syscall.O_RDONLY|syscall.O_NOCTTY, 0666)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("handleReadFile")
	}

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

	logrus.WithFields(logrus.Fields{
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
	mutex.Lock()
	teleinfo.Entries[trame.etiquette] = trame.data
	mutex.Unlock()
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

// get values
func (teleinfo *Teleinfo) GetEntries(entries map[string]string) map[string]string {
	mutex.Lock()
	for key, value := range teleinfo.Entries {
		entries[key] = value
	}
	mutex.Unlock()
	return entries
}

// get values
func (teleinfo *Teleinfo) Get(key string) string {
	var value string
	mutex.Lock()
	i, ex := teleinfo.Entries[key]
	if ex {
		value = i
	}
	mutex.Unlock()
	return value
}

// initialize this module
func (that *Teleinfo) init() {
	// add map
	that.Entries = make(map[string]string)

	// start worker
	go handleReadFile(getTeleinfoFile())
	go worker(that)

	// log information
	log.Default.Info("teleinfo", log.Fields{
		"teleinfoFile": getTeleinfoFile(),
	})
}
