package cron

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"

	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	log "github.com/yroffin/jarvis-go-ext/logger"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/mongodb"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/razberry"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/teleinfo"
	mgo "gopkg.in/mgo.v2"
)

// AdvertiseJob
type AdvertiseJob struct {
}

// Run
func (job *AdvertiseJob) Run() {

	// define default value for this connector
	m := &types.Connector{
		Name:       viper.GetString("jarvis.module.name"),
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

	// check for s
	if errs != nil {
		log.Default.Error("cron", log.Fields{
			"s": errs,
		})
		return
	}

	// check for s
	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Default.Error("cron", log.Fields{
			"body":   string(b),
			"status": resp.Status,
		})
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
		log.Default.Error("teleinfo", log.Fields{
			"data": &CollectTeleinfoResource{Base: base, Timestamp: time.Now()},
			"":     err,
		})
	}
}

type CollectRazberryJob struct {
	mgo      *mongodb.MongoDriver
	col      *mgo.Collection
	razberry *razberry.Razberry
	devices  []string
}

// CollectTeleinfoResource : CollectTeleinfoResource resource struct
type CollectRazberryResource struct {
	Timestamp time.Time
	Name      string
	Device    map[string]interface{}
}

// Run the job CollectTeleinfoJob
func (job *CollectRazberryJob) Run() {
	/**
	 * store data
	 */
	for index := 0; index < len(job.devices); index++ {
		var dev, err = job.razberry.DeviceById(job.devices[index])
		if err == nil {
			job.col.Insert(&CollectRazberryResource{Name: job.devices[index], Device: dev, Timestamp: time.Now()})
		} else {
			log.Default.Error("razberry", log.Fields{
				"data": &CollectRazberryResource{Name: job.devices[index], Device: dev, Timestamp: time.Now()},
				"":     err,
			})
		}
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
	// advertise
	if viper.GetString("jarvis.option.advertise") == "true" {
		// first call
		var job = new(AdvertiseJob)
		job.Run()
		// init cron
		c := cron.New()
		c.AddJob(viper.GetString("jarvis.option.advertise.cron"), job)
		c.Start()
		logrus.WithFields(logrus.Fields{
			"cron": viper.GetString("jarvis.option.advertise.cron"),
		}).Info("advertise")
	}

	// teleinfo
	if viper.GetString("jarvis.option.teleinfo.active") == "true" {
		// first call
		var job = new(CollectTeleinfoJob)
		/**
		* store mongo session
		 */
		job.mgo = mongodb.GetInstance()
		job.col = job.mgo.GetCollection("collect", "teleinfo")
		job.teleinfo = teleinfo.GetInstance()
		job.Run()

		// init cron
		c := cron.New()
		c.AddJob(viper.GetString("jarvis.option.teleinfo.cron"), job)
		c.Start()
		logrus.WithFields(logrus.Fields{
			"cron": viper.GetString("jarvis.option.teleinfo.cron"),
		}).Info("teleinfo")
	}

	// teleinfo
	if viper.GetString("jarvis.option.razberry.active") == "true" {
		// first call
		var job = new(CollectRazberryJob)
		/**
		* store mongo session
		 */
		job.mgo = mongodb.GetInstance()
		job.col = job.mgo.GetCollection("collect", "razberry")
		job.razberry = razberry.GetInstance()
		job.devices = strings.Split(viper.GetString("jarvis.option.razberry.devices"), ",")
		job.Run()

		// init cron
		c := cron.New()
		c.AddJob(viper.GetString("jarvis.option.razberry.cron"), job)
		c.Start()
		logrus.WithFields(logrus.Fields{
			"cron": viper.GetString("jarvis.option.razberry.cron"),
		}).Info("razberry")
	}
}
