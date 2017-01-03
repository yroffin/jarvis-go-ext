package cron

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/mongodb"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/teleinfo"
	mgo "gopkg.in/mgo.v2"
)

type AdvertiseJob struct {
}

// AdvertiseJob : Run
func (job *AdvertiseJob) Run() {

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
		logrus.WithFields(logrus.Fields{
			"errors": errs,
		}).Error("CRON")
		return
	}

	if b, err := ioutil.ReadAll(resp.Body); err == nil {
		logrus.WithFields(logrus.Fields{
			"body":   string(b),
			"status": resp.Status,
		}).Debug("CRON")
	} else {
		logrus.WithFields(logrus.Fields{
			"body":   string(b),
			"status": resp.Status,
		}).Warn("WARN")
	}
}

type CollectTeleinfoJob struct {
	mgo      *mongodb.MongoDriver
	col      *mgo.Collection
	teleinfo *teleinfo.Teleinfo
}

// CollectTeleinfoResource : CollectTeleinfoResource resource struct
type CollectTeleinfoResource struct {
	Timestamp time.Time
	Base      int
}

// Run the job CollectTeleinfoJob
func (job *CollectTeleinfoJob) Run() {
	/**
	 * store data
	 */
	var base, err = strconv.Atoi(job.teleinfo.Get("BASE"))
	if err == nil {
		job.col.Insert(&CollectTeleinfoResource{Base: base, Timestamp: time.Now()})
	} else {
		logrus.WithFields(logrus.Fields{
			"data":  &CollectTeleinfoResource{Base: base, Timestamp: time.Now()},
			"error": err,
		}).Error("Teleinfo")
	}
}

// CronDriver : cron driver instance
type CronDriver struct {
}

var instance *CronDriver
var once sync.Once

// GetInstance : singleton instance
func GetInstance() *CronDriver {
	once.Do(func() {
		instance = new(CronDriver)
		instance.initAdvertise()
	})
	return instance
}

// InitAdvertise : init cron service
func (cronDriver *CronDriver) initAdvertise() {
	var advertise = viper.GetString("jarvis.option.advertise")

	if advertise != "" {
		// first call
		var job = new(AdvertiseJob)
		job.Run()
		// init cron
		c := cron.New()
		c.AddJob(advertise, job)
		c.Start()
	}

	/**
	 * teleinfo
	 */
	var teleinfoString = viper.GetString("jarvis.option.teleinfo.collect")
	if teleinfoString != "" {
		// first call
		var job = new(CollectTeleinfoJob)
		/**
		 * store mongo session
		 */
		job.mgo = mongodb.GetInstance()
		job.col = job.mgo.GetCollection("collect", "teleinfo")
		job.teleinfo = teleinfo.GetInstance()
		job.Run()

		logrus.WithFields(logrus.Fields{
			"job": job,
		}).Info("Teleinfo")

		// init cron
		c := cron.New()
		c.AddJob(teleinfoString, job)
		c.Start()
	}

	logrus.WithFields(logrus.Fields{
		"advertise": advertise,
		"teleinfo":  teleinfoString,
	}).Info("CronDriver")
}
