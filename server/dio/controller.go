package dio

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/mfrc522"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/wiringpi"
)

// HandlePost : handler for post
func HandlePostDio(c echo.Context) error {
	var m *types.DioResource
	c.Bind(&m)

	logger.NewLogger().WithFields(logrus.Fields{
		"pin":         m.Pin,
		"sender":      m.Sender,
		"interruptor": m.Interuptor,
		"on":          m.On,
	}).Info("DIO")

	if m.On == true {
		wiringpi.On(m.Pin, m.Sender, m.Interuptor)
	} else {
		wiringpi.Off(m.Pin, m.Sender, m.Interuptor)
	}

	return c.JSON(http.StatusOK, m)
}

// HandlePost : handler for post
func HandlePostSpi(c echo.Context) error {
	var m *types.DioResource
	c.Bind(&m)

	logger.NewLogger().WithFields(logrus.Fields{
		"pin":         m.Pin,
		"sender":      m.Sender,
		"interruptor": m.Interuptor,
		"on":          m.On,
	}).Info("SPI")

	return c.JSON(http.StatusOK, m)
}

// HandlePost : handler for post
func HandlePostMfrc522(c echo.Context) error {
	var m *types.Mfrc522Resource
	c.Bind(&m)

	logger.NewLogger().WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522")

	return c.JSON(http.StatusOK, m)
}

// HandlePostMfrc522DumpClassic1K : handler for post
func HandlePostMfrc522DumpClassic1K(c echo.Context) error {
	var m *types.Mfrc522Resource
	c.Bind(&m)

	logger.NewLogger().WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522/dump")

	// request for status
	var instance = mfrc522.GetInstance()
	m.Status, m.Data = instance.Request(mfrc522.PICC_REQIDL)
	m.Len = len(m.Data)
	if m.Status != 0 {
		return c.JSON(http.StatusBadRequest, m)
	}

	logger.NewLogger().WithFields(logrus.Fields{}).Info("mfrc522/dump/request")

	// request for uid
	m.Status, m.Data = instance.Anticoll()
	m.Len = len(m.Data)
	if m.Status != 0 {
		return c.JSON(http.StatusBadRequest, m)
	}

	m.Uid[0] = m.Data[0]
	m.Uid[1] = m.Data[1]
	m.Uid[2] = m.Data[2]
	m.Uid[3] = m.Data[3]
	m.Uid[4] = m.Data[4]

	logger.NewLogger().WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522/dump/anticoll")

	// Select tag
	instance.SelectTag(m.Uid)

	// DumpClassic1K
	instance.DumpClassic1K(m.Key, m.Uid)

	return c.JSON(http.StatusOK, m)
}

// HandlePostMfrc522Request : handler for post
func HandlePostMfrc522Request(c echo.Context) error {
	var m *types.Mfrc522Resource
	c.Bind(&m)

	logger.NewLogger().WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522/request")

	var instance = mfrc522.GetInstance()
	m.Status, m.Data = instance.Request(mfrc522.PICC_REQIDL)
	m.Len = len(m.Data)

	return c.JSON(http.StatusOK, m)
}

// HandlePostMfrc522AntiColl : handler for post
func HandlePostMfrc522AntiColl(c echo.Context) error {
	var m *types.Mfrc522Resource
	c.Bind(&m)

	logger.NewLogger().WithFields(logrus.Fields{
		"key": m.Key,
		"uid": m.Uid,
	}).Info("mfrc522/anticoll")

	var instance = mfrc522.GetInstance()
	m.Status, m.Data = instance.Anticoll()
	m.Len = len(m.Data)
	if m.Status == 0 {
		m.Uid[0] = m.Data[0]
		m.Uid[1] = m.Data[1]
		m.Uid[2] = m.Data[2]
		m.Uid[3] = m.Data[3]
	}

	return c.JSON(http.StatusOK, m)
}
