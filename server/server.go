package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	ctrlDio "github.com/yroffin/jarvis-go-ext/server/dio"
	"github.com/yroffin/jarvis-go-ext/server/utils/cron"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
	"github.com/yroffin/jarvis-go-ext/server/utils/native"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

// Start : start the jarvis server
func Start() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api")
	{ // routes for /api
		dio := api.Group("/dio")
		{ // routes for /api/dio
			dio.Post("", ctrlDio.HandlePost)
		}
	}

	// init wiringPi library
	native.InitWiringPi()

	// init cron
	cron.Init("@every 60s")

	// get prot from config
	intf := viper.GetString("jarvis.module.interface")
	port := viper.GetString("jarvis.module.port")

	logger.NewLogger().WithFields(log.Fields{
		"interface": intf,
		"port":      port,
	}).Info("DIO")

	e.Run(standard.New(intf + ":" + port))
}
