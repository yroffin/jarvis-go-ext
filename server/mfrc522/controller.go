package mfrc522

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/mfrc522"
)

// HandlePost : handler for post
func HandlePostMfrc522(c echo.Context) error {
	var m *types.Mfrc522Resource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522")

	return c.JSON(http.StatusOK, m)
}

// HandlePostMfrc522DumpClassic1K : handler for post
func HandlePostMfrc522DumpClassic1K(c echo.Context) error {
	var m *types.Mfrc522DumpResource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522/dump")

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
func HandlePostMfrc522WriteClassic1K(c echo.Context) error {
	var m *types.Mfrc522WriteResource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522/dump")

	return c.JSON(http.StatusOK, m)
}

// HandlePostMfrc522Request : handler for post
func HandlePostMfrc522Request(c echo.Context) error {
	var m *types.Mfrc522DumpResource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"key": m.Key,
	}).Info("HandlePostMfrc522Request")

	return c.JSON(http.StatusOK, m)
}

// HandlePostMfrc522AntiColl : handler for post
func HandlePostMfrc522AntiColl(c echo.Context) error {
	var m *types.Mfrc522Resource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("HandlePostMfrc522AntiColl")

	return c.JSON(http.StatusOK, m)
}
