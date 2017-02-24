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

package cron_service

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/yroffin/jarvis-go-ext/logger"
	"github.com/yroffin/jarvis-go-ext/server/service/mongodb_service"
	"github.com/yroffin/jarvis-go-ext/server/service/razberry_service"
	"github.com/yroffin/jarvis-go-ext/server/service/teleinfo_service"
	"github.com/yroffin/jarvis-go-ext/server/types"
	mgo "gopkg.in/mgo.v2"
)

// CronService service descriptor
type CronService struct {
}

var instance *CronService
var once sync.Once
var mutex = &sync.Mutex{}

// CronService singleton instance
func Service() *CronService {
	once.Do(func() {
		instance = new(CronService)
		instance.init()
	})
	return instance
}

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
		logger.Default.Error("cron", logger.Fields{
			"s": errs,
		})
		return
	}

	// check for s
	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		logger.Default.Error("cron", logger.Fields{
			"body":   string(b),
			"status": resp.Status,
		})
	}
}

// CollectTeleinfoJob collect job
type CollectTeleinfoJob struct {
	mgo      *mongodb_service.MongoService
	col      *mgo.Collection
	teleinfo *teleinfo_service.TeleinfoService
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
		logger.Default.Error("teleinfo", logger.Fields{
			"data": &CollectTeleinfoResource{Base: base, Timestamp: time.Now()},
			"":     err,
		})
	}
}

// CollectRazberryJob collect job
type CollectRazberryJob struct {
	mgo      *mongodb_service.MongoService
	col      *mgo.Collection
	razberry *razberry_service.RazberryService
	devices  []string
}

// CollectRazberryResource collect resource
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
			logger.Default.Error("razberry", logger.Fields{
				"data": &CollectRazberryResource{Name: job.devices[index], Device: dev, Timestamp: time.Now()},
				"":     err,
			})
		}
	}
}

// InitAdvertise : init cron service
func (that *CronService) init() {
	// advertise
	if viper.GetString("jarvis.option.advertise") == "true" {
		// first call
		var job = new(AdvertiseJob)
		job.Run()
		// init cron
		c := cron.New()
		c.AddJob(viper.GetString("jarvis.option.advertise.cron"), job)
		c.Start()
		logger.Default.Info("advertise", logger.Fields{
			"cron": viper.GetString("jarvis.option.advertise.cron"),
		})
	}

	// teleinfo
	if viper.GetString("jarvis.option.teleinfo.active") == "true" {
		// first call
		var job = new(CollectTeleinfoJob)
		/**
		* store mongo session
		 */
		job.mgo = mongodb_service.Service()
		job.col = job.mgo.GetCollection("collect", "teleinfo")
		job.teleinfo = teleinfo_service.Service()
		job.Run()

		// init cron
		c := cron.New()
		c.AddJob(viper.GetString("jarvis.option.teleinfo.cron"), job)
		c.Start()
		logger.Default.Info("teleinfo", logger.Fields{
			"cron": viper.GetString("jarvis.option.advertise.cron"),
		})
	}

	// teleinfo
	if viper.GetString("jarvis.option.razberry.active") == "true" {
		// first call
		var job = new(CollectRazberryJob)
		/**
		* store mongo session
		 */
		job.mgo = mongodb_service.Service()
		job.col = job.mgo.GetCollection("collect", "razberry")
		job.razberry = razberry_service.Service()
		job.devices = strings.Split(viper.GetString("jarvis.option.razberry.devices"), ",")
		job.Run()

		// init cron
		c := cron.New()
		c.AddJob(viper.GetString("jarvis.option.razberry.cron"), job)
		c.Start()
		logger.Default.Info("razberry", logger.Fields{
			"cron": viper.GetString("jarvis.option.advertise.cron"),
		})
	}
}
