package cron

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
)

// handlerAdvertise : handlerAdvertise for connector register
func handlerAdvertise() {

	// define default value for this connector
	m := &types.Connector{
		Name:       "go-dio",
		Icon:       "settings_input_antenna",
		Adress:     "http://" + viper.GetString("jarvis.module.interface") + ":" + viper.GetString("jarvis.module.port"),
		IsRenderer: true,
		IsSensor:   false,
		CanAnswer:  false,
	}

	mJSON, _ := json.Marshal(m)

	request := gorequest.New().Timeout(2 * time.Second)
	resp, _, errs := request.
		Post(viper.GetString("jarvis.server.url") + "/api/connectors/*?task=register").
		Send(string(mJSON)).
		End()
	if errs != nil {
		logger.NewLogger().WithFields(logrus.Fields{
			"errors": errs,
		}).Error("CRON")
		return
	}

	if b, err := ioutil.ReadAll(resp.Body); err == nil {
		logger.NewLogger().WithFields(logrus.Fields{
			"body":   string(b),
			"status": resp.Status,
		}).Debug("CRON")
	} else {
		logger.NewLogger().WithFields(logrus.Fields{
			"body":   string(b),
			"status": resp.Status,
		}).Warn("WARN")
	}
}

// InitAdvertise : init cron service
func InitAdvertise(cr string) int {
	// first call
	handlerAdvertise()
	// init cron
	c := cron.New()
	c.AddFunc(cr, handlerAdvertise)
	c.Start()
	return 0
}
