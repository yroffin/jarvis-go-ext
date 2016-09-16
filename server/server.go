package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	ctrlDio "github.com/yroffin/jarvis-go-ext/server/dio"
	"github.com/yroffin/jarvis-go-ext/server/utils/cron"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/mfrc522"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/wiringpi"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

// Start : start the jarvis server
func Start() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if viper.GetString("jarvis.option.mfrc522") == "true" {
		// init mfrc522 singleton
		mfrc522.GetInstance()
		logger.NewLogger().WithFields(logrus.Fields{
			"actived": "true",
		}).Info("mfrc522")
	}

	if viper.GetString("jarvis.option.wiringpi") == "true" {
		// init wiringPi library
		wiringpi.Init()
		logger.NewLogger().WithFields(logrus.Fields{
			"actived": "true",
		}).Info("wiringpi")
	}

	if viper.GetString("jarvis.option.advertise") == "true" {
		// init cron
		cron.Init("@every 60s")
		logger.NewLogger().WithFields(logrus.Fields{
			"actived": "true",
		}).Info("cron")
	}

	api := e.Group("/api")
	{ // routes for /api
		dio := api.Group("/dio")
		{ // routes for /api/dio
			dio.Post("", ctrlDio.HandlePostDio)
		}
		spi := api.Group("/spi")
		{ // routes for /api/spi
			spi.Post("", ctrlDio.HandlePostSpi)
		}
	}

	// get prot from config
	intf := viper.GetString("jarvis.module.interface")
	port := viper.GetString("jarvis.module.port")

	logger.NewLogger().WithFields(logrus.Fields{
		"interface": intf,
		"port":      port,
	}).Info("module")

	e.Run(standard.New(intf + ":" + port))
}
