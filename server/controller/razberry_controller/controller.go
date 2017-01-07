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

package razberry_controller

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/service/razberry_service"
	"github.com/yroffin/jarvis-go-ext/server/types"
)

// Get handler for get
func Get(c echo.Context) error {
	var m *types.RazberryResource
	m = new(types.RazberryResource)
	c.Bind(&m)

	if c.Param("id") == "" {
		var body, err = razberry_service.Service().Devices()
		if err != nil {
			return c.JSON(http.StatusBadRequest, m)
		}
		return c.JSON(http.StatusOK, body)
	} else {
		var body, err = razberry_service.Service().DeviceById(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, m)
		}
		return c.JSON(http.StatusOK, body)
	}
}
