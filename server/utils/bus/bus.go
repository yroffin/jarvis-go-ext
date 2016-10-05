package bus

import (
	"fmt"
	"time"

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
		}
	}
}

func Start() {
	go produce()
	go consume()
}
