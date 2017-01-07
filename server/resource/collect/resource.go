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

package collect

import (
	"net/http"
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/service/mongodb"
	"github.com/yroffin/jarvis-go-ext/server/types"
)

// MongodbService struct
type CollectService struct {
	mongoMiddleware *mongodb.MongoDriver
}

var instance *CollectService
var once sync.Once

// GetInstance : singleton instance
func GetInstance() *CollectService {
	once.Do(func() {
		instance = new(CollectService)
		instance.init()
	})
	return instance
}

// Get return all local collect
func (that *CollectService) GetAll(c echo.Context) error {
	var m *types.CollectResource
	m = new(types.CollectResource)
	c.Bind(&m)

	// retrieve all collections stored in "collect" database
	var names, err = that.mongoMiddleware.GetCollections("collect")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// reduce this collection by removing system.indexes
	var reduce = make([]string, len(names))
	var counter = 0
	for index := 1; index < len(names); index++ {
		if names[index] != "system.indexes" {
			reduce[counter] = names[index]
			counter++
		}
	}

	var result = make([]string, counter)
	copy(result, reduce[0:counter])

	// result
	m.Collections = result

	return c.JSON(http.StatusOK, m)
}

// Get all collection elements
func (that *CollectService) Get(c echo.Context) error {
	var m *types.CollectResource
	m = new(types.CollectResource)
	c.Bind(&m)

	if c.Param("id") == "" {
		// result
		var names, err = that.findAll()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		m.Collections = names
	} else {
		// result
		if c.QueryParam("orderby") == "" {
			var names, err = that.get(c.Param("id"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			m.Data = names
		} else {
			var names, err = that.getAndSort(c.Param("id"), c.QueryParam("orderby"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			m.Data = names
		}
	}

	return c.JSON(http.StatusOK, m)
}

// findAll find all collection name in database
func (that *CollectService) findAll() ([]string, error) {
	// retrieve all collections stored in "collect" database
	var names, err = that.mongoMiddleware.GetCollections("collect")
	if err != nil {
		return nil, err
	}

	// reduce this collection by removing system.indexes
	var reduce = make([]string, len(names))
	var counter = 0
	for index := 1; index < len(names); index++ {
		if names[index] != "system.indexes" {
			reduce[counter] = names[index]
			counter++
		}
	}

	var result = make([]string, counter)
	copy(result, reduce[0:counter])
	return result, nil
}

// findAll find all collection name in database
func (that *CollectService) get(name string) ([]bson.M, error) {
	// retrieve all collections stored in "collect" database
	tuples := []bson.M{}
	that.mongoMiddleware.GetCollection("collect", name).Find(bson.M{}).All(&tuples)
	return tuples, nil
}

// findAll find all collection name in database
func (that *CollectService) getAndSort(name string, sort string) ([]bson.M, error) {
	// retrieve all collections stored in "collect" database
	tuples := []bson.M{}
	that.mongoMiddleware.GetCollection("collect", name).Find(bson.M{}).Sort(sort).All(&tuples)
	return tuples, nil
}

// init initialize service
func (that *CollectService) init() error {
	that.mongoMiddleware = mongodb.GetInstance()
	return nil
}
