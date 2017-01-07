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

package mongodb

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
)

// MongoDriver : mongo driver instance
type MongoDriver struct {
	session *mgo.Session
}

var instance *MongoDriver
var once sync.Once

// GetInstance : singleton instance
func GetInstance() *MongoDriver {
	once.Do(func() {
		instance = new(MongoDriver)
		instance.init()
	})
	return instance
}

// get collections
func (mongoDriver *MongoDriver) GetCollections(db string) ([]string, error) {
	return mongoDriver.session.DB(db).CollectionNames()
}

// get collection
func (mongoDriver *MongoDriver) GetCollection(db string, col string) *mgo.Collection {
	return mongoDriver.session.DB(db).C(col)
}

// store element
func (mongoDriver *MongoDriver) StoreData(db string, col string, data interface{}) interface{} {
	return mongoDriver.Store(mongoDriver.GetCollection(db, col), data)
}

// store element
func (mongoDriver *MongoDriver) Store(col *mgo.Collection, data interface{}) interface{} {
	var err = col.Insert(&data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"data":  data,
			"error": err,
		}).Error("MongoDriver")
	}
	return data
}

// Close session
func (mongoDriver *MongoDriver) Close() {
	defer mongoDriver.session.Close()
}

// initialize this module
func (mongoDriver *MongoDriver) init() {
	var host = viper.GetString("jarvis.option.mongodb")

	logrus.WithFields(logrus.Fields{
		"host": host,
	}).Info("MongoDriver")

	// get mongo session
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	mongoDriver.session = session

	var info, _ = session.BuildInfo()

	logrus.WithFields(logrus.Fields{
		"host":    host,
		"version": info.Version,
		"sys":     info.SysInfo,
	}).Info("MongoDriver")
}
