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

package teleinfo

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/teleinfo"
)

// HandleGet : handler for get
func HandleGetTeleinfo(c echo.Context) error {
	var m *types.TeleinfoResource
	m = new(types.TeleinfoResource)
	c.Bind(&m)

	var instance = teleinfo.GetInstance()
	m.Entries = make(map[string]string)
	instance.GetEntries(m.Entries)

	return c.JSON(http.StatusOK, m)
}
