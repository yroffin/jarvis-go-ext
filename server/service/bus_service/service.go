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

package bus_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"

	"github.com/yroffin/jarvis-go-ext/logger"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/mfrc522"
)

// BusService service descriptor
type BusService struct {
}

var instance *BusService
var once sync.Once
var mutex = &sync.Mutex{}

// BusService singleton instance
func Service() *BusService {
	once.Do(func() {
		instance = new(BusService)
		instance.init()
	})
	return instance
}

var msgs = make(chan *types.MessageResource)

func produce() {
	var instance = mfrc522.GetInstance()

	for {
		var msg = new(types.MessageResource)
		// dump nfc tag
		var result, tagType, uuid = instance.WaitForTag()
		if result != nil {
			logger.Default.Warn("module", logger.Fields{
				"status":  result,
				"message": "Unable to detect tag",
			})
			msg.TagType = "None"
			msg.TagUuid = "0x0000000000"
		} else {
			msg.TagType = tagType
			msg.TagUuid = fmt.Sprintf("0x%02x%02x%02x%02x%02x", uuid[0], uuid[1], uuid[2], uuid[3], uuid[4])
		}
		msgs <- msg

		time.Sleep(2 * time.Second)
	}
}

// consume : consume message
func consume() {
	var last types.MessageResource

	for {
		msg := <-msgs

		if msg.TagUuid != last.TagUuid {
			// New Tag detection
			logger.Default.Info("consume", logger.Fields{
				"message": msg,
			})
			last.TagType = msg.TagType
			last.TagUuid = msg.TagUuid

			if last.TagType != "None" {
				mJSON, _ := json.Marshal(last)

				request := gorequest.New().Timeout(2 * time.Second)
				resp, _, errs := request.
					Patch(viper.GetString("jarvis.server.url") + "/api/triggers/nfc:" + last.TagUuid).
					Send(string(mJSON)).
					End()

				if errs != nil {
					logger.Default.Error("nfc", logger.Fields{
						"errors": errs,
					})
				} else {
					if b, err := ioutil.ReadAll(resp.Body); err == nil {
						logger.Default.Info("nfc", logger.Fields{
							"url":    viper.GetString("jarvis.server.url") + "/api/triggers/nfc:" + last.TagUuid,
							"body":   string(b),
							"status": resp.Status,
						})
					} else {
						logger.Default.Error("nfc", logger.Fields{
							"url":    viper.GetString("jarvis.server.url") + "/api/triggers/nfc:" + last.TagUuid,
							"body":   string(b),
							"status": resp.Status,
						})
					}
				}
			}
		}
	}
}

// initialize this module
func (that *BusService) init() {
	go produce()
	go consume()

	// log information
	logger.Default.Info("bus", logger.Fields{})
}
