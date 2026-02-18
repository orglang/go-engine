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
	bindingErr := c.Bind(&dto)
	if bindingErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindingErr
	}
	ctx := c.Request().Context()
	h.log.Log(ctx, lf.LevelTrace, "posting started", slog.Any("dto", dto))
	validationErr := dto.Validate()
	if validationErr != nil {
		h.log.Error("validation failed", slog.Any("dto", dto))
		return validationErr
	}
	spec, conversionErr := MsgToDefSpec(dto)
	if conversionErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return conversionErr
	}
	snap, creationErr := h.api.Create(spec)
	if creationErr != nil {
		return creationErr
	}
	h.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("ref", snap.DescRef))
	return c.JSON(http.StatusCreated, MsgFromDefSnap(snap))
}
