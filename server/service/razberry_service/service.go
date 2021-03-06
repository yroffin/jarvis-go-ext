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

package razberry_service

import (
	"errors"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"github.com/yroffin/jarvis-go-ext/logger"
)

// RazberryService service descriptor
type RazberryService struct {
	Url  string
	Auth string
}

var instance *RazberryService
var once sync.Once
var mutex = &sync.Mutex{}

// RazberryService singleton instance
func Service() *RazberryService {
	once.Do(func() {
		instance = new(RazberryService)
		instance.init()
	})
	return instance
}

// Devices get device by id
func (that *RazberryService) Devices() (map[string]interface{}, error) {
	return that.get("/devices")
}

// Devices get device by id
func (that *RazberryService) DeviceById(id string) (map[string]interface{}, error) {
	var index = strings.Split(id, "_")
	var indexes = strings.Split(index[2], "-")
	return that.get("/devices[" + indexes[0] + "].instances[" + indexes[1] + "].commandClasses[" + indexes[2] + "].data[" + indexes[3] + "]")
}

// Get
func (that *RazberryService) get(uri string) (map[string]interface{}, error) {

	logger.Default.Debug("razberryService", logger.Fields{
		"uri": that.Url + uri,
	})

	request := gorequest.New().Timeout(2 * time.Second)

	resp, _, errs := request.
		Get(that.Url+uri).
		Set("Authorization", that.Auth).
		End()

	// check for errors
	if errs != nil {
		logger.Default.Error("razberry", logger.Fields{
			"errors": errs,
		})
		return nil, errors.New("http: while trying to connect")
	}

	// check for errors
	if b, err := ioutil.ReadAll(resp.Body); err != nil || resp.StatusCode != 200 {
		logger.Default.Error("razberry", logger.Fields{
			"body":   string(b),
			"status": resp.Status,
		})
		return nil, errors.New("http: while decoding body")
	} else {
		var body map[string]interface{}
		json.Unmarshal(b, &body)
		return body, nil
	}
}

// initialize this module
func (that *RazberryService) init() {
	that.Url = viper.GetString("jarvis.option.razberry.url")
	that.Auth = viper.GetString("jarvis.option.razberry.auth")

	// log information
	logger.Default.Info("razberry", logger.Fields{
		"url": that.Url,
	})
}
