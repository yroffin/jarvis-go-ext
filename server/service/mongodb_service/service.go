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
	"fmt"
	"sync"

	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
)

// MongoService : mongo driver instance
type MongoService struct {
	session *mgo.Session
}

var instance *MongoService
var once sync.Once

// Service : singleton instance
func Service() *MongoService {
	once.Do(func() {
		instance = new(MongoService)
		instance.init()
	})
	return instance
}

// Insert : get collections
func (MongoService *MongoService) Insert(col *mgo.Collection, docs interface{}) error {
	var err error
	err = nil
	// try one insert
	err = col.Insert(docs)
	if err != nil {
		// refresh session if needed
		MongoService.session.Refresh()
		// retry another time
		err = col.Insert(docs)
	}
	return err
}

// GetCollections : get collections
func (MongoService *MongoService) GetCollections(db string) ([]string, error) {
	// refresh session if needed
	MongoService.session.Refresh()
	return MongoService.session.DB(db).CollectionNames()
}

// GetCollection : get collection
func (MongoService *MongoService) GetCollection(db string, col string) *mgo.Collection {
	// refresh session if needed
	MongoService.session.Refresh()
	return MongoService.session.DB(db).C(col)
}

// StoreData : store element
func (MongoService *MongoService) StoreData(db string, col string, data interface{}) interface{} {
	return MongoService.Store(MongoService.GetCollection(db, col), data)
}

// store element
func (MongoService *MongoService) Store(col *mgo.Collection, data interface{}) interface{} {
	var err = col.Insert(&data)
	if err != nil {
		fmt.Printf("[ERROR] while backup data %s/%s\n", data, err)
	}
	return data
}

// Close session
func (MongoService *MongoService) Close() {
	defer MongoService.session.Close()
}

// init initialize this module
func (that *MongoService) init() {
	var host = viper.GetString("jarvis.option.mongodb")

	fmt.Printf("[INFO] init mongodb %s\n", host)

	// get mongo session
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	that.session = session

	var info, _ = session.BuildInfo()

	fmt.Printf("[INFO] init mongodb %s/%s/%s\n", host, info.Version, info.SysInfo)
}
