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
