package dio

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
	"github.com/yroffin/jarvis-go-ext/server/utils/native"
)

// HandlePost : handler for post
func HandlePost(c echo.Context) error {
	var m *types.DioResource
	c.Bind(&m)

	logger.NewLogger().WithFields(log.Fields{
		"pin":         m.Pin,
		"sender":      m.Sender,
		"interruptor": m.Interuptor,
		"on":          m.On,
	}).Info("DIO")

	if m.On == "true" {
		native.DioOn(m.Pin, m.Sender, m.Interuptor)
	} else {
		native.DioOff(m.Pin, m.Sender, m.Interuptor)
	}

	c.JSON(http.StatusOK, m)
}
