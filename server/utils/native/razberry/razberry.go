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

package razberry

import (
	"errors"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	log "github.com/yroffin/jarvis-go-ext/logger"
)

// Razberry : instance Razberry device struct
type Razberry struct {
	Url  string
	Auth string
}

var instance *Razberry
var once sync.Once
var mutex = &sync.Mutex{}

// GetInstance : singleton instance
func GetInstance() *Razberry {
	once.Do(func() {
		instance = new(Razberry)
		instance.init()
	})
	return instance
}

// Devices get device by id
func (that *Razberry) Devices() (map[string]interface{}, error) {
	return that.get("/devices")
}

// Devices get device by id
func (that *Razberry) DeviceById(id string) (map[string]interface{}, error) {
	var index = strings.Split(id, "_")
	var indexes = strings.Split(index[2], "-")
	return that.get("/devices[" + indexes[0] + "].instances[" + indexes[1] + "].commandClasses[" + indexes[2] + "].data[" + indexes[3] + "]")
}

// Get
func (that *Razberry) get(uri string) (map[string]interface{}, error) {

	logrus.WithFields(logrus.Fields{
		"uri": that.Url + uri,
	}).Debug("razberry")

	request := gorequest.New().Timeout(2 * time.Second)

	resp, _, errs := request.
		Get(that.Url+uri).
		Set("Authorization", that.Auth).
		End()

	// check for errors
	if errs != nil {
		log.Default.Error("razberry", log.Fields{
			"errors": errs,
		})
		logrus.WithFields(logrus.Fields{
			"errors": errs,
		}).Error("razberry")
		return nil, errors.New("http: while trying to connect")
	}

	// check for errors
	if b, err := ioutil.ReadAll(resp.Body); err != nil || resp.StatusCode != 200 {
		log.Default.Error("razberry", log.Fields{
			"body":   string(b),
			"status": resp.Status,
		})
		logrus.WithFields(logrus.Fields{
			"body":   string(b),
			"status": resp.Status,
		}).Error("razberry")
		return nil, errors.New("http: while decoding body")
	} else {
		var body map[string]interface{}
		json.Unmarshal(b, &body)
		return body, nil
	}
}

// initialize this module
func (that *Razberry) init() {
	that.Url = viper.GetString("jarvis.option.razberry.url")
	that.Auth = viper.GetString("jarvis.option.razberry.auth")

	// log information
	log.Default.Info("razberry", log.Fields{})
}
