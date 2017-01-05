package dio

import (
	"net/http"

	"github.com/labstack/echo"
	log "github.com/yroffin/jarvis-go-ext/logger"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/wiringpi"
)

// Post handle post on dio resource
func Post(c echo.Context) error {
	var m *types.DioResource
	c.Bind(&m)

	log.Default.Info("dio", log.Fields{
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
