package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	ctrlDio "github.com/yroffin/jarvis-go-ext/server/dio"
	ctrlMfrc522 "github.com/yroffin/jarvis-go-ext/server/mfrc522"
	ctrlTeleinfo "github.com/yroffin/jarvis-go-ext/server/teleinfo"
	bus "github.com/yroffin/jarvis-go-ext/server/utils/bus"
	"github.com/yroffin/jarvis-go-ext/server/utils/cron"
	"github.com/yroffin/jarvis-go-ext/server/utils/mongodb"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/mfrc522"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/teleinfo"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/wiringpi"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	log "github.com/yroffin/jarvis-go-ext/logger"
)

// Start : start the jarvis server
func Start() {
	e := echo.New()
	e.Use(middleware.Recover())

	// setup mongodb
	if viper.GetString("jarvis.option.mongodb") != "" {
		// init MongoDb driver
		mongodb.GetInstance()
		logrus.WithFields(logrus.Fields{
			"active": "true",
		}).Info("mongodb")
	}

	// init mongodb logs
	e.Use(log.GetInstance().GetMiddleware())

	// setup wiring pi
	if viper.GetString("jarvis.option.wiringpi") == "true" {
		// init wiringPi library
		wiringpi.GetInstance()
		log.Default.Info("wiringpi", log.Fields{
			"active": "true",
		})
	}

	// setup mfrc522
	if viper.GetString("jarvis.option.mfrc522") == "true" {
		// init mfrc522 singleton
		mfrc522.GetInstance()
		log.Default.Info("mfrc522", log.Fields{
			"active": "true",
		})
	}

	// setup teleinfo
	if viper.GetString("jarvis.option.teleinfo") == "true" {
		// init teleinfo singleton
		teleinfo.GetInstance()
		log.Default.Info("teleinfo", log.Fields{
			"active": "true",
		})
	}

	// setup cron manager
	cron.GetInstance()

	log.Default.Info("cron", log.Fields{
		"active": "true",
	})

	api := e.Group("/api")
	{ // routes for /api
		dio := api.Group("/dio")
		{ // routes for /api/dio
			dio.Post("", ctrlDio.HandlePostDio)
		}
		teleinfo := api.Group("/teleinfo")
		{ // routes for /api/dio
			teleinfo.Get("", ctrlTeleinfo.HandleGetTeleinfo)
		}
		mfrc522 := api.Group("/mfrc522")
		{ // routes for /api/mfrc522
			mfrc522.Post("", ctrlMfrc522.HandlePostMfrc522)
			mfrc522Anticoll := mfrc522.Group("/anticoll")
			{ // routes for /api/mfrc522/anticoll
				mfrc522Anticoll.Post("", ctrlMfrc522.HandlePostMfrc522AntiColl)
			}
			mfrc522Request := mfrc522.Group("/request")
			{ // routes for /api/mfrc522/request
				mfrc522Request.Post("", ctrlMfrc522.HandlePostMfrc522Request)
			}
			mfrc522DumpClassic1K := mfrc522.Group("/dump")
			{ // routes for /api/mfrc522/dump
				mfrc522DumpClassic1K.Post("", ctrlMfrc522.HandlePostMfrc522DumpClassic1K)
			}
		}
	}

	// get prot from config
	intf := viper.GetString("jarvis.module.interface")
	port := viper.GetString("jarvis.module.port")

	logrus.WithFields(logrus.Fields{
		"interface": intf,
		"port":      port,
	}).Info("module")

	if viper.GetString("jarvis.option.nfctag") == "true" {
		// start nfc capture
		bus.Start()
		log.Default.Info("nfctag", log.Fields{
			"active": "true",
		})
	}

	e.Run(standard.New(intf + ":" + port))
}
