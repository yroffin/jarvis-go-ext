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

	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/logger"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/mfrc522"
)

// Post : handler for post
func Post(c echo.Context) error {
	var m *types.Mfrc522Resource
	c.Bind(&m)

	logger.Default.Info("mfrc522/post", logger.Fields{
		"key": m.Key,
		"uid": m.Uid,
	})

	return c.JSON(http.StatusOK, m)
}

// PostDumpClassic1K : handler for post
func PostDumpClassic1K(c echo.Context) error {
	var m *types.Mfrc522DumpResource
	c.Bind(&m)

	logger.Default.Info("mfrc522/dump", logger.Fields{
		"key": m.Key,
		"uid": m.Uid,
	})

	// get native service (low level)
	var instance = mfrc522.GetInstance()

	// dump nfc tag
	var result, tagType, data, _ = instance.DumpClassic1K(m.Key[0:len(m.Key)])
	if result != nil {
		logger.Default.Error("mfrc522/dump", logger.Fields{
			"status": result,
			"error":  "Unable to detect tag",
		})
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

// PostWriteClassic1K : handler for post
func PostWriteClassic1K(c echo.Context) error {
	var m *types.Mfrc522WriteResource
	c.Bind(&m)

	logger.Default.Info("mfrc522/write", logger.Fields{
		"key": m.Key,
		"uid": m.Uid,
	})

	return c.JSON(http.StatusOK, m)
}

// PostRequest : handler for post
func PostRequest(c echo.Context) error {
	var m *types.Mfrc522DumpResource
	c.Bind(&m)

	logger.Default.Info("mfrc522/request", logger.Fields{
		"key": m.Key,
		"uid": m.Uid,
	})

	return c.JSON(http.StatusOK, m)
}

// PostAntiColl : handler for post
func PostAntiColl(c echo.Context) error {
	var m *types.Mfrc522Resource
	c.Bind(&m)

	logger.Default.Info("mfrc522/anticoll", logger.Fields{
		"key": m.Key,
		"uid": m.Uid,
	})

	return c.JSON(http.StatusOK, m)
}
