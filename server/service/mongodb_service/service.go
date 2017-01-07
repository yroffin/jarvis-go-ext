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

package mongodb_service

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
)

// MongoService : mongo driver instance
type MongoService struct {
	session *mgo.Session
}

var instance *MongoService
var once sync.Once

// GetInstance : singleton instance
func Service() *MongoService {
	once.Do(func() {
		instance = new(MongoService)
		instance.init()
	})
	return instance
}

// get collections
func (MongoService *MongoService) GetCollections(db string) ([]string, error) {
	return MongoService.session.DB(db).CollectionNames()
}

// get collection
func (MongoService *MongoService) GetCollection(db string, col string) *mgo.Collection {
	return MongoService.session.DB(db).C(col)
}

// store element
func (MongoService *MongoService) StoreData(db string, col string, data interface{}) interface{} {
	return MongoService.Store(MongoService.GetCollection(db, col), data)
}

// store element
func (MongoService *MongoService) Store(col *mgo.Collection, data interface{}) interface{} {
	var err = col.Insert(&data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"data":  data,
			"error": err,
		}).Error("MongoService")
	}
	return data
}

// Close session
func (MongoService *MongoService) Close() {
	defer MongoService.session.Close()
}

// initialize this module
func (MongoService *MongoService) init() {
	var host = viper.GetString("jarvis.option.mongodb")

	logrus.WithFields(logrus.Fields{
		"host": host,
	}).Info("MongoService")

	// get mongo session
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	MongoService.session = session

	var info, _ = session.BuildInfo()

	logrus.WithFields(logrus.Fields{
		"host":    host,
		"version": info.Version,
		"sys":     info.SysInfo,
	}).Info("MongoService")
}
