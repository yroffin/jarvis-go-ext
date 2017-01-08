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

package dio_controller

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/logger"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/wiringpi"
)

// Post handle post on dio resource
func Post(c echo.Context) error {
	var m *types.DioResource
	c.Bind(&m)

	logrus.WithFields(logrus.Fields{
		"c":    c,
		"bind": m,
	}).Info("dio")

	logger.Default.Info("dio", logger.Fields{
		"pin":         m.Pin,
		"sender":      m.Sender,
		"interruptor": m.Interuptor,
		"on":          m.On,
	})

	if m.On == true {
		wiringpi.On(m.Pin, m.Sender, m.Interuptor)
	} else {
		wiringpi.Off(m.Pin, m.Sender, m.Interuptor)
	}

	return c.JSON(http.StatusOK, m)
}
