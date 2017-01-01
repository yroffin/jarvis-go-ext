package teleinfo

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/types"
)

// HandleGet : handler for get
func HandleGetTeleinfo(c echo.Context) error {
	var m *types.TeleinfoResource
	c.Bind(&m)

	return c.JSON(http.StatusOK, m)
}
