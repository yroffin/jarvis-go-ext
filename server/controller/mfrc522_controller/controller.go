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

package mfrc522_controller

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/mfrc522"
)

// HandlePost : handler for post
func Post(c echo.Context) error {
	var m *types.Mfrc522Resource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522")

	return c.JSON(http.StatusOK, m)
}

// HandlePostMfrc522DumpClassic1K : handler for post
func PostDumpClassic1K(c echo.Context) error {
	var m *types.Mfrc522DumpResource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522/dump")

	// get native service (low level)
	var instance = mfrc522.GetInstance()

	// dump nfc tag
	var result, tagType, data, _ = instance.DumpClassic1K(m.Key[0:len(m.Key)])
	if result != nil {
		logrus.WithFields(logrus.Fields{
			"status": result,
		}).Error("Unable to detect tag")
	}

	// fill result
	m.TagType = tagType
	for index := 0; index < len(data); index++ {
		m.Uid[index] = data[index]
	}

	// for sector := 0; sector < len(sectors); sector++ {
	// for value := 0; value < len(sectors[sector].Values); value++ {
	// m.Sectors[sector].Values[value] = sectors[sector].Values[value]
	// }
	// }

	return c.JSON(http.StatusOK, m)
}

// HandlePostMfrc522WriteClassic1K : handler for post
func PostWriteClassic1K(c echo.Context) error {
	var m *types.Mfrc522WriteResource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522/dump")

	return c.JSON(http.StatusOK, m)
}

// HandlePostMfrc522Request : handler for post
func PostRequest(c echo.Context) error {
	var m *types.Mfrc522DumpResource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"key": m.Key,
	}).Info("HandlePostMfrc522Request")

	return c.JSON(http.StatusOK, m)
}

// HandlePostMfrc522AntiColl : handler for post
func PostAntiColl(c echo.Context) error {
	var m *types.Mfrc522Resource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("HandlePostMfrc522AntiColl")

	return c.JSON(http.StatusOK, m)
}
