package bus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"

	"github.com/Sirupsen/logrus"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/mfrc522"
)

var msgs = make(chan *types.MessageResource)

func produce() {
	var instance = mfrc522.GetInstance()

	for {
		var msg = new(types.MessageResource)
		// dump nfc tag
		var result, tagType, uuid = instance.WaitForTag()
		if result != nil {
			logger.NewLogger().WithFields(logrus.Fields{
				"status": result,
			}).Warn("Unable to detect tag")
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
			logger.NewLogger().WithFields(logrus.Fields{
				"message": msg,
			}).Info("consume")
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
					logger.NewLogger().WithFields(logrus.Fields{
						"errors": errs,
					}).Error("nfc")
				} else {
					if b, err := ioutil.ReadAll(resp.Body); err == nil {
						logger.NewLogger().WithFields(logrus.Fields{
							"url":    viper.GetString("jarvis.server.url") + "/api/triggers/nfc:" + last.TagUuid,
							"body":   string(b),
							"status": resp.Status,
						}).Info("nfc")
					} else {
						logger.NewLogger().WithFields(logrus.Fields{
							"url":    viper.GetString("jarvis.server.url") + "/api/triggers/nfc:" + last.TagUuid,
							"body":   string(b),
							"status": resp.Status,
						}).Error("nfc")
					}
				}
			}
		}
	}
}

func Start() {
	go produce()
	go consume()
}
