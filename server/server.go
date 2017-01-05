package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	dioCtrl "github.com/yroffin/jarvis-go-ext/server/resource/dio"
	mfrc522Ctrl "github.com/yroffin/jarvis-go-ext/server/resource/mfrc522"
	razberryCtrl "github.com/yroffin/jarvis-go-ext/server/resource/razberry"
	teleinfoCtrl "github.com/yroffin/jarvis-go-ext/server/resource/teleinfo"
	bus "github.com/yroffin/jarvis-go-ext/server/utils/bus"
	"github.com/yroffin/jarvis-go-ext/server/utils/cron"
	"github.com/yroffin/jarvis-go-ext/server/utils/mongodb"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/mfrc522"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/razberry"
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
		// init mfrc522 service
		mfrc522.GetInstance()
		log.Default.Info("mfrc522", log.Fields{
			"active": "true",
		})
	}

	// setup teleinfo
	if viper.GetString("jarvis.option.teleinfo") == "true" {
		// init teleinfo service
		teleinfo.GetInstance()
		log.Default.Info("teleinfo", log.Fields{
			"active": "true",
		})
	}

	// setup razberry
	if viper.GetString("jarvis.option.razberry") == "true" {
		// init razberry service
		razberry.GetInstance()
		log.Default.Info("razberry", log.Fields{
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
		if viper.GetString("jarvis.option.dio") == "true" {
			logrus.WithFields(logrus.Fields{
				"interface": "dio",
			}).Info("module")
			dioGroup := api.Group("/dio")
			{ // routes for /api/dio
				dioGroup.Post("", dioCtrl.HandlePostDio)
			}
		}
		if viper.GetString("jarvis.option.teleinfo") == "true" {
			logrus.WithFields(logrus.Fields{
				"interface": "teleinfo",
			}).Info("module")
			teleinfoGroup := api.Group("/teleinfo")
			{ // routes for /api/teleinfo
				teleinfoGroup.Get("", teleinfoCtrl.HandleGetTeleinfo)
			}
		}
		if viper.GetString("jarvis.option.razberry") == "true" {
			logrus.WithFields(logrus.Fields{
				"interface": "razberry",
			}).Info("module")
			razberryGroup := api.Group("/razberry")
			{ // routes for /api/razberry
				razberryGroup.Get("", razberryCtrl.Get)
				razberryGroup.Get("/:id", razberryCtrl.Get)
			}
		}
		if viper.GetString("jarvis.option.mfrc522") == "true" {
			logrus.WithFields(logrus.Fields{
				"interface": "mfrc522",
			}).Info("module")
			mfrc522Group := api.Group("/mfrc522")
			{ // routes for /api/mfrc522
				mfrc522Group.Post("", mfrc522Ctrl.HandlePostMfrc522)
				mfrc522Anticoll := mfrc522Group.Group("/anticoll")
				{ // routes for /api/mfrc522/anticoll
					mfrc522Anticoll.Post("", mfrc522Ctrl.HandlePostMfrc522AntiColl)
				}
				mfrc522Request := mfrc522Group.Group("/request")
				{ // routes for /api/mfrc522/request
					mfrc522Request.Post("", mfrc522Ctrl.HandlePostMfrc522Request)
				}
				mfrc522DumpClassic1K := mfrc522Group.Group("/dump")
				{ // routes for /api/mfrc522/dump
					mfrc522DumpClassic1K.Post("", mfrc522Ctrl.HandlePostMfrc522DumpClassic1K)
				}
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
