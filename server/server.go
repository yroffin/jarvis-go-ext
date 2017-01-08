package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yroffin/jarvis-go-ext/server/controller/collect_controller"
	"github.com/yroffin/jarvis-go-ext/server/controller/dio_controller"
	"github.com/yroffin/jarvis-go-ext/server/controller/mfrc522_controller"
	"github.com/yroffin/jarvis-go-ext/server/controller/razberry_controller"
	"github.com/yroffin/jarvis-go-ext/server/controller/teleinfo_controller"
	"github.com/yroffin/jarvis-go-ext/server/service/bus_service"
	"github.com/yroffin/jarvis-go-ext/server/service/cron_service"
	"github.com/yroffin/jarvis-go-ext/server/service/mongodb_service"
	"github.com/yroffin/jarvis-go-ext/server/service/razberry_service"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/mfrc522"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/wiringpi"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	log "github.com/yroffin/jarvis-go-ext/logger"
	"github.com/yroffin/jarvis-go-ext/server/service/teleinfo_service"
)

// Start : start the jarvis server
func Start() {
	e := echo.New()
	e.Use(middleware.Recover())

	// setup mongodb
	if viper.GetString("jarvis.option.mongodb") != "" {
		// init MongoDb driver
		mongodb_service.Service()
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
		teleinfo_service.Service()
		log.Default.Info("teleinfo", log.Fields{
			"active": "true",
		})
	}

	// setup razberry
	if viper.GetString("jarvis.option.razberry") == "true" {
		// init razberry service
		razberry_service.Service()
		log.Default.Info("razberry", log.Fields{
			"active": "true",
		})
	}

	// setup cron manager
	cron_service.Service()
	log.Default.Info("cron", log.Fields{
		"active": "true",
	})

	logrus.WithFields(logrus.Fields{
		"interface": "api",
	}).Info("module")

	api := e.Group("/api")
	{ // routes for /api
		collectGroup := api.Group("/collect")
		{ // routes for /api/collect
			collectGroup.Post("/:id", collect_controller.Post)
			collectGroup.Get("", collect_controller.Get)
		}
		if viper.GetString("jarvis.option.dio") == "true" {
			logrus.WithFields(logrus.Fields{
				"interface": "dio",
			}).Info("module")
			dioGroup := api.Group("/dio")
			{ // routes for /api/dio
				dioGroup.Post("", dio_controller.Post)
			}
		}
		if viper.GetString("jarvis.option.teleinfo") == "true" {
			logrus.WithFields(logrus.Fields{
				"interface": "teleinfo",
			}).Info("module")
			teleinfoGroup := api.Group("/teleinfo")
			{ // routes for /api/teleinfo
				teleinfoGroup.Get("", teleinfo_controller.Get)
			}
		}
		if viper.GetString("jarvis.option.razberry") == "true" {
			logrus.WithFields(logrus.Fields{
				"interface": "razberry",
			}).Info("module")
			api.Get("/razberry/:id", razberry_controller.Get)
			api.Get("/razberry", razberry_controller.Get)
		}
		if viper.GetString("jarvis.option.mfrc522") == "true" {
			logrus.WithFields(logrus.Fields{
				"interface": "mfrc522",
			}).Info("module")
			mfrc522Group := api.Group("/mfrc522")
			{ // routes for /api/mfrc522
				mfrc522Group.Post("", mfrc522_controller.Post)
				mfrc522Anticoll := mfrc522Group.Group("/anticoll")
				{ // routes for /api/mfrc522/anticoll
					mfrc522Anticoll.Post("", mfrc522_controller.PostAntiColl)
				}
				mfrc522Request := mfrc522Group.Group("/request")
				{ // routes for /api/mfrc522/request
					mfrc522Request.Post("", mfrc522_controller.PostRequest)
				}
				mfrc522DumpClassic1K := mfrc522Group.Group("/dump")
				{ // routes for /api/mfrc522/dump
					mfrc522DumpClassic1K.Post("", mfrc522_controller.PostDumpClassic1K)
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
		bus_service.Service()
		log.Default.Info("nfctag", log.Fields{
			"active": "true",
		})
	}

	e.Run(standard.New(intf + ":" + port))
}
