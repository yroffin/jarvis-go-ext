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

package collect_controller

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/service/mongodb_service"
	"github.com/yroffin/jarvis-go-ext/server/types"
)

// GetAll return all local collect
func GetAll(c echo.Context) error {
	var m *types.CollectResource
	m = new(types.CollectResource)
	c.Bind(&m)

	// retrieve all collections stored in "collect" database
	var names, err = mongodb_service.Service().GetCollections("collect")
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
func Get(c echo.Context) error {
	var m *types.CollectResource
	m = new(types.CollectResource)
	c.Bind(&m)

	if c.Param("id") == "" {
		// result
		var names, err = findAll()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		m.Collections = names
	} else {
		// result
		if c.QueryParam("orderby") == "" {
			var names, err = get(c.Param("id"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			m.Data = names
		} else {
			var names, err = getAndSort(c.Param("id"), c.QueryParam("orderby"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			m.Data = names
		}
	}

	return c.JSON(http.StatusOK, m)
}

// findAll find all collection name in database
func findAll() ([]string, error) {
	// retrieve all collections stored in "collect" database
	var names, err = mongodb_service.Service().GetCollections("collect")
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
func get(name string) ([]bson.M, error) {
	// retrieve all collections stored in "collect" database
	tuples := []bson.M{}
	mongodb_service.Service().GetCollection("collect", name).Find(bson.M{}).All(&tuples)
	return tuples, nil
}

// findAll find all collection name in database
func getAndSort(name string, sort string) ([]bson.M, error) {
	// retrieve all collections stored in "collect" database
	tuples := []bson.M{}
	mongodb_service.Service().GetCollection("collect", name).Find(bson.M{}).Sort(sort).All(&tuples)
	return tuples, nil
}
