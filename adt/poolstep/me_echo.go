package poolstep

import (
	"log/slog"
	"reflect"

	"github.com/labstack/echo/v4"
)

// Server-side primary adapter
type echoController struct {
	api API
	log *slog.Logger
}

func newEchoController(api API, log *slog.Logger) *echoController {
	name := slog.String("name", reflect.TypeFor[echoController]().Name())
	return &echoController{api, log.With(name)}
}

func cfgEchoController(e *echo.Echo, h *echoController) error {
	return nil
}
