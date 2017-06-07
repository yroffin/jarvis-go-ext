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

package mqtt_service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/yroffin/jarvis-go-ext/logger"
)

// MqttService service descriptor
type MqttService struct {
	handle mqtt.Client
}

var instance *MqttService
var once sync.Once
var mutex = &sync.Mutex{}

// Service singleton instance
func Service() *MqttService {
	once.Do(func() {
		instance = new(MqttService)
		instance.init()
	})
	return instance
}

// Publish text
func (that *MqttService) Publish(topic string, payload string) *MqttService {
	return that.PublishRaw(topic, payload, 0, false)
}

// Publish text
func (that *MqttService) PublishRaw(topic string, payload string, qos byte, retained bool) *MqttService {
	var token = that.handle.Publish(topic, qos, retained, payload)
	token.Wait()
	return that
}

// PublishData
func (that *MqttService) PublishData(topic string, jsonData interface{}) *MqttService {
	b, err := json.MarshalIndent(jsonData, "", " ")
	if err != nil {
		fmt.Println("error:", err)
	}
	that.Publish(topic, string(b))
	return that
}

// internal handler
var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	logger.Default.Info("mqtt", logger.Fields{
		"topic":   msg.Topic(),
		"payload": msg.Payload(),
	})
}

// init initialize this module
func (that *MqttService) init() {
	// mqtt client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(viper.GetString("jarvis.mqtt.url"))
	opts.SetClientID("gojarvis")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	that.handle = mqtt.NewClient(opts)
	if token := that.handle.Connect(); token.Wait() && token.Error() != nil {
		// log error and panic
		logger.Default.Error("mqtt", logger.Fields{
			"url":         viper.GetString("jarvis.mqtt.url"),
			"id":          "gojarvis",
			"token/error": token.Error(),
		})
		panic(token.Error())
	} else {
		// log information
		logger.Default.Info("mqtt", logger.Fields{
			"url": viper.GetString("jarvis.mqtt.url"),
			"id":  "gojarvis",
		})
	}
}
