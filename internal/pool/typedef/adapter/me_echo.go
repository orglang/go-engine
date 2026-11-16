package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdktypedef "github.com/orglang/go-sdk/pool/typedef"

	"orglang/go-engine/internal/lib/lf"

	"orglang/go-engine/internal/adt/typesem"

	typedef "orglang/go-engine/internal/pool/typedef/core"
)

// Controller is the inbound HTTP port implemented by the Echo adapter.
type Controller interface {
	PostSpec(echo.Context) error
}

// RegisterRoutes attaches the pool type-definition endpoint to Echo.
func RegisterRoutes(e *echo.Echo, h Controller) error {
	e.POST("/api/v1/xacts/defs", h.PostSpec)
	return nil
}

// Server-side primary adapter
type controllerEcho struct {
	api typedef.API
	log *slog.Logger
}

func newControllerEcho(a typedef.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{a, l.With(name)}
}

func cfgEchoController(e *echo.Echo, h *controllerEcho) error {
	return RegisterRoutes(e, h)
}

func (h *controllerEcho) PostSpec(c echo.Context) error {
	var dto sdktypedef.DefSpec
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	ctx := c.Request().Context()
	h.log.Log(ctx, lf.LevelTrace, "posting started", slog.Any("dto", dto))
	validateErr := dto.Validate()
	if validateErr != nil {
		h.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
	}
	spec, convErr := MsgToDefSpec(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	ref, createErr := h.api.Create(spec)
	if createErr != nil {
		return createErr
	}
	h.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("ref", ref))
	return c.JSON(http.StatusCreated, typesem.MsgFromRef(ref))
}
