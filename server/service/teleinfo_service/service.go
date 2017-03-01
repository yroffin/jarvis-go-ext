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

package teleinfo_service

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/viper"
	log "github.com/yroffin/jarvis-go-ext/logger"
)

// On linux
// apply: stty 1200 cs7 evenp cstopb -igncr -inlcr -brkint -icrnl -opost -isig -icanon -iexten -F /dev/ttyUSB0
// to configure teleinfo tty
// Cf. https://hallard.me/gestion-de-la-teleinfo-avec-un-raspberry-pi-et-une-carte-arduipi/

// Teleinfo : instance Teleinfo device struct
type TeleinfoService struct {
	Entries map[string]string
}

// Trame EDF
type TeleinfoTrame struct {
	etiquette string // ETIQUETTE (4 à 8 caractères)
	data      string // DATA (1 à 12 caractères)
	hecksum   string // CHECKSUM (caractère de contrôle : Somme de caractère)
}

var instance *TeleinfoService
var once sync.Once
var mutex = &sync.Mutex{}

// GetInstance : singleton instance
func Service() *TeleinfoService {
	once.Do(func() {
		instance = new(TeleinfoService)
		instance.init()
	})
	return instance
}

// declare canal
var canal = make(chan byte, 1024)

// handleReadFile : read file
func handleReadFile(device string) error {

	s, err := os.OpenFile(device, syscall.O_RDONLY, 0666)

	if err != nil {
		log.Default.Error("teleinfo", log.Fields{
			"error": err,
		})
	}

	log.Default.Info("teleinfo", log.Fields{
		"device": device,
	})

	buffer := make([]byte, 4096)
	reader := bufio.NewReader(s)
	var millis = time.Millisecond
	millis = 1

	// Receive reply
	for {
		if _, err := reader.Read(buffer); err != nil {
			// sleep while no bytes
			// to avoid system flood read
			log.Default.Error("teleinfo", log.Fields{
				"Error": err,
			})
			fmt.Printf("e")
			time.Sleep(1000 * time.Millisecond)
		} else {
			// dispatch io
			for i := 0; i < len(buffer); i++ {
				canal <- buffer[i]
			}
			if len(buffer) < 3000 {
				fmt.Printf(">")
				millis = millis + 10
			} else {
				fmt.Printf("<")
				millis = millis - 10
				if millis == 0 {
					millis = 0
				}
			}
			time.Sleep(millis * time.Millisecond)
		}
	}

	log.Default.Info("teleinfo", log.Fields{
		"status": "done",
	})

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
func (that *TeleinfoService) submit(trame TeleinfoTrame) {
	mutex.Lock()
	that.Entries[trame.etiquette] = trame.data
	mutex.Unlock()
}

// single trame
func (that *TeleinfoService) handleTrame(trame string) {
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
	that.submit(send)
}

// all trames detection
func (that *TeleinfoService) handleTrames(trame string) {
	var send string
	for i := 0; i < len(trame); i++ {
		switch {
		case trame[i] == 0x0A:
			send = ""
			continue
		case trame[i] == 0x0D:
			that.handleTrame(send)
			continue
		default:
			send += string([]byte{trame[i]})
		}
	}
}

// worker to consume file
func worker(that *TeleinfoService) {
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
			that.handleTrames(trame)
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

// GetEntries load entries
func (that *TeleinfoService) GetEntries(entries map[string]string) map[string]string {
	mutex.Lock()
	for key, value := range that.Entries {
		entries[key] = value
	}
	mutex.Unlock()
	return entries
}

// Get get values
func (that *TeleinfoService) Get(key string) string {
	var value string
	mutex.Lock()
	i, ex := that.Entries[key]
	if ex {
		value = i
	}
	mutex.Unlock()
	return value
}

// init initialize this module
func (that *TeleinfoService) init() {
	// add map
	that.Entries = make(map[string]string)

	// start worker
	go handleReadFile(viper.GetString("jarvis.option.teleinfo.file"))
	go worker(that)

	// log information
	log.Default.Info("teleinfo", log.Fields{
		"teleinfoFile": viper.GetString("jarvis.option.teleinfo.file"),
	})
}
