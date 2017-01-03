package logger

import (
	"strconv"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/utils/mongodb"
	mgo "gopkg.in/mgo.v2"
)

// LoggerTools : instance logger
type LoggerTools struct {
	mgo         *mongodb.MongoDriver
	collections map[string]*mgo.Collection
}

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// LogResource : data to log
type LogResource struct {
	Timestamp time.Time
	Fields    Fields
}

// Default logger
var Default *LoggerTools
var instance *LoggerTools
var once sync.Once
var mutex = &sync.Mutex{}

// GetMiddleware : singleton instance
func (that *LoggerTools) GetMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			that.Info("access", Fields{
				"remote_ip":     req.RealIP(),
				"host":          req.Host(),
				"uri":           req.URI(),
				"method":        req.Method(),
				"path":          req.URL().Path(),
				"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
				"latency_human": stop.Sub(start).String(),
				"bytes_in":      req.Header().Get(echo.HeaderContentLength),
				"bytes_out":     strconv.FormatInt(res.Size(), 10),
			})

			return nil
		}
	}
}

// GetInstance : singleton instance
func GetInstance() *LoggerTools {
	once.Do(func() {
		instance = new(LoggerTools)
		Default = instance
		instance.init()
	})
	return instance
}

// Info : log info data
func (that *LoggerTools) Info(category string, fields Fields) {
	mutex.Lock()
	// add collection to map il not exist
	if _, ok := that.collections[category]; !ok {
		that.collections[category] = that.mgo.GetCollection("logger", category)
	}
	mutex.Unlock()
	that.collections[category].Insert(&LogResource{Fields: fields, Timestamp: time.Now()})
}

// Init : Init
func (that *LoggerTools) init() {
	logrus.WithFields(logrus.Fields{}).Info("LoggerTools")

	that.mgo = mongodb.GetInstance()
	that.collections = make(map[string]*mgo.Collection)
}
