package xactdef

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"github.com/orglang/go-sdk/adt/xactdef"

	"orglang/go-engine/lib/lf"
)

// Server-side primary adapter
type echoController struct {
	api API
	log *slog.Logger
}

func newEchoController(a API, l *slog.Logger) *echoController {
	name := slog.String("name", reflect.TypeFor[echoController]().Name())
	return &echoController{a, l.With(name)}
}

func cfgEchoController(e *echo.Echo, h *echoController) error {
	e.POST("/api/v1/xacts/defs", h.PostSpec)
	return nil
}

func (h *echoController) PostSpec(c echo.Context) error {
	var dto xactdef.DefSpec
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
	spec, convertErr := MsgToDefSpec(dto)
	if convertErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convertErr
	}
	snap, createErr := h.api.Create(spec)
	if createErr != nil {
		return createErr
	}
	h.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("ref", snap.DescRef))
	return c.JSON(http.StatusCreated, MsgFromDefSnap(snap))
}
