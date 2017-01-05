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

package razberry

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/razberry"
)

// HandleGet : handler for get
func Get(c echo.Context) error {
	var m *types.RazberryResource
	m = new(types.RazberryResource)
	c.Bind(&m)

	var razberry = razberry.GetInstance()

	if c.Value("id") == "" {
		var body, err = razberry.Devices()
		if err != nil {
			return c.JSON(http.StatusBadRequest, m)
		}
		return c.JSON(http.StatusOK, body)
	} else {
		var body, err = razberry.DeviceById(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, m)
		}
		return c.JSON(http.StatusOK, body)
	}
}
