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

package logger

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/service/mongodb_service"
	mgo "gopkg.in/mgo.v2"
)

// LoggerTools : instance logger
type LoggerTools struct {
	mgo         *mongodb_service.MongoService
	collections map[string]*mgo.Collection
}

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// LogResource : data to log
type LogResource struct {
	Timestamp time.Time
	Fields    Fields
}

// Default logger
var Default *LoggerTools
var instance *LoggerTools
var once sync.Once
var mutex = &sync.Mutex{}

// GetMiddleware : singleton instance
func (that *LoggerTools) GetMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			that.Info("access", Fields{
				"remote_ip":     req.RealIP(),
				"host":          req.Host(),
				"uri":           req.URI(),
				"method":        req.Method(),
				"path":          req.URL().Path(),
				"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
				"latency_human": stop.Sub(start).String(),
				"bytes_in":      req.Header().Get(echo.HeaderContentLength),
				"bytes_out":     strconv.FormatInt(res.Size(), 10),
			})

			return nil
		}
	}
}

// GetInstance : singleton instance
func GetInstance() *LoggerTools {
	once.Do(func() {
		instance = new(LoggerTools)
		Default = instance
		instance.init()
	})
	return instance
}

// Info : log info data
func (that *LoggerTools) InfoLog(category string, fields Fields) {
	fields["category"] = category
	fields["timestamp"] = time.Now()
	jsonString, _ := json.Marshal(fields)
	that.InfoSyslog(string(jsonString))
}

// Info : log info data
func (that *LoggerTools) Info(category string, fields Fields) {
	mutex.Lock()
	// add collection to map il not exist
	if _, ok := that.collections[category]; !ok {
		that.collections[category] = that.mgo.GetCollection("logger", category)
	}
	mutex.Unlock()
	fields["Level"] = "INFO"
	that.collections[category].Insert(&LogResource{Fields: fields, Timestamp: time.Now()})
	that.InfoLog(category, fields)
}

// Debug : log error data
func (that *LoggerTools) Debug(category string, fields Fields) {
	mutex.Lock()
	// add collection to map il not exist
	if _, ok := that.collections[category]; !ok {
		that.collections[category] = that.mgo.GetCollection("logger", category)
	}
	mutex.Unlock()
	fields["Level"] = "DEBUG"
	that.collections[category].Insert(&LogResource{Fields: fields, Timestamp: time.Now()})
}

// Error : log error data
func (that *LoggerTools) Error(category string, fields Fields) {
	mutex.Lock()
	// add collection to map il not exist
	if _, ok := that.collections[category]; !ok {
		that.collections[category] = that.mgo.GetCollection("logger", category)
	}
	mutex.Unlock()
	fields["Level"] = "ERROR"
	that.collections[category].Insert(&LogResource{Fields: fields, Timestamp: time.Now()})
	that.InfoLog(category, fields)
}

// Warn : log error data
func (that *LoggerTools) Warn(category string, fields Fields) {
	mutex.Lock()
	// add collection to map il not exist
	if _, ok := that.collections[category]; !ok {
		that.collections[category] = that.mgo.GetCollection("logger", category)
	}
	mutex.Unlock()
	fields["Level"] = "WARN"
	that.collections[category].Insert(&LogResource{Fields: fields, Timestamp: time.Now()})
	that.InfoLog(category, fields)
}

// Init : Init
func (that *LoggerTools) init() {
	that.mgo = mongodb_service.Service()
	that.collections = make(map[string]*mgo.Collection)
	that.initSyslog()
}
