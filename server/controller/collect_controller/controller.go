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
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Sirupsen/logrus"
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

	// X-Total-Count
	c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(m.Collections)))
	return c.JSON(http.StatusOK, m)
}

// parse any data
func parse(field map[string]interface{}, key string, value interface{}) error {
	// check for key not null
	if value.(map[string]interface{})[key] != nil {
		// no extract object {"format": "value""}
		for keyToDecode, valueToDecode := range value.(map[string]interface{})[key].(map[string]interface{}) {
			var decoded interface{}
			var err error
			if keyToDecode == "RFC3339" {
				decoded, err = time.Parse(time.RFC3339, valueToDecode.(string))
			}
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err.Error(),
				}).Error("collect")
				return err
			}
			field[key] = decoded
		}
	}
	return nil
}

// Post all collection elements
func Post(c echo.Context) error {
	if c.QueryParam("operation") == "find" {
		return find(c)
	}
	if c.QueryParam("operation") == "pipe" {
		return pipeOperation(c)
	}
	return nil
}

// aggregate all collection elements
func pipeOperation(c echo.Context) error {
	var m *types.PostCollectResource
	c.Bind(&m)

	var data, err = pipe(c.Param("id"), m.Pipes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// X-Total-Count
	c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(data)))
	return c.JSON(http.StatusOK, data)
}

// find all collection elements
func find(c echo.Context) error {
	var m *types.PostCollectResource
	c.Bind(&m)

	var q = map[string]interface{}{}
	for key, value := range m.Find {
		var field = map[string]interface{}{}
		var err error
		// $gt handler
		err = parse(field, "$gt", value)
		if err != nil {
			return c.JSON(http.StatusBadRequest, bson.M{"error": err.Error(), "detail": err})
		}
		// $lt handler
		err = parse(field, "$lt", value)
		if err != nil {
			return c.JSON(http.StatusBadRequest, bson.M{"error": err.Error(), "detail": err})
		}
		q[key] = field
	}

	var r = new(types.CollectResource)

	// result
	if m.OrderBy == nil {
		var data, err = get(c.Param("id"), q)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		r.Data = data
	} else {
		var data, err = getAndSort(c.Param("id"), q, m.OrderBy)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		r.Data = data
	}

	// X-Total-Count
	c.Response().Header().Set("X-Total-Count", strconv.Itoa(len(r.Data)))
	return c.JSON(http.StatusOK, r)
}

// Get all collection elements
func Get(c echo.Context) error {
	// result
	var details, err = findAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var body = new(types.GetCollectResource)
	body.Collections = details

	return c.JSON(http.StatusOK, body)
}

// findAll find all collection name in database
func findAll() ([]types.GetCollectResourceDetail, error) {
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

	// result contains all collections
	var result = make([]string, counter)
	copy(result, reduce[0:counter])

	// convert it to GetCollectResourceDetail
	var arr = make([]types.GetCollectResourceDetail, len(result))
	for index := 0; index < len(result); index++ {
		arr[index].Name = result[index]
		// find one sample entity
		tuples := []bson.M{}
		mongodb_service.Service().
			GetCollection("collect", arr[index].Name).
			Find(&bson.M{}).
			Sort("-$natural").
			Limit(1).
			All(&tuples)
		if len(tuples) > 0 {
			arr[index].Entity = tuples[0]
		}
	}

	return arr, nil
}

// findAll find all collection name in database
func get(name string, m bson.M) ([]bson.M, error) {
	// retrieve all collections stored in "collect" database
	tuples := []bson.M{}
	mongodb_service.Service().GetCollection("collect", name).Find(m).All(&tuples)
	return tuples, nil
}

// findAll find all collection name in database
func getAndSort(name string, m bson.M, sort []string) ([]bson.M, error) {
	// retrieve all collections stored in "collect" database
	tuples := []bson.M{}
	mongodb_service.Service().GetCollection("collect", name).Find(m).Sort(sort[0]).All(&tuples)
	return tuples, nil
}

// dump this map
func apply(m []bson.M) {
	for k, v := range m {
		// Wrap the original in a reflect.Value
		translateRecursive(0, m, k, v)
	}
}

// translateRecursive apply IISODate conversion
func translateRecursive(level int, holder interface{}, key interface{}, value interface{}) {
	switch value.(type) {
	case string:
		// convert strig when type conversion is needed
		if strings.HasPrefix(value.(string), "ISODate(") {

			extract := strings.Replace(value.(string), "ISODate(", "", 1)
			extractFinal := strings.Replace(extract, ")", "", 1)

			dateFormat, err := time.Parse(time.RFC3339, extractFinal)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err.Error(),
				}).Error("collect")
			} else {
				map[string]interface{}(holder.(map[string]interface{}))[string(key.(string))] = dateFormat
			}
		}
		break
	case bson.M:
		for k, v := range value.(bson.M) {
			// Wrap bson.M
			translateRecursive(level+1, value, k, v)
		}
		break
	case map[string]interface{}:
		for k, v := range value.(map[string]interface{}) {
			// Wrap map[string]interface{}
			translateRecursive(level+1, value, k, v)
		}
		break
	default:
	}
}

// pipe aggregate functionor map reduce
func pipe(name string, m []bson.M) ([]bson.M, error) {
	// apply pipes on collections stored in "collect" database
	tuples := []bson.M{}
	// transform data
	apply(m)
	mongodb_service.Service().GetCollection("collect", name).Pipe(m).All(&tuples)
	return tuples, nil
}
