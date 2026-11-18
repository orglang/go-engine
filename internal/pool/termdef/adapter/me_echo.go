package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/app/web/api"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

// Controller is the inbound HTTP port implemented by the Echo adapter.
type Controller interface {
	PostSpec(echo.Context) error
}

// RegisterRoutes attaches the pool term-definition endpoint to Echo.
func RegisterRoutes(e *echo.Echo, h Controller) error {
	e.POST("/api/v1/pools/declarations", h.PostSpec)
	return nil
}

// Server-side primary adapter
type controllerEcho struct {
	api termdef.API
	log *slog.Logger
}

func newControllerEcho(a termdef.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{a, l.With(name)}
}

func cfgEchoController(e *echo.Echo, h *controllerEcho) error {
	return RegisterRoutes(e, h)
}

func (h *controllerEcho) PostSpec(c echo.Context) error {
	var dto api.PoolDecSpec
	if err := c.Bind(&dto); err != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return err
	}
	spec, err := DecSpecFromGenerated(dto)
	if err != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return err
	}
	ref, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, termsem.MsgFromRef(ref))
}
