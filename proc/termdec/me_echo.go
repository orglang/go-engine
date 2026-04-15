package termdec

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdk "github.com/orglang/go-sdk/adt/descsem"
	"github.com/orglang/go-sdk/adt/procdec"

	"orglang/go-engine/adt/semtype"
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
	e.POST("/api/v1/procs/decs", h.PostSpec)
	e.GET("/api/v1/procs/decs/:id", h.GetSnap)
	return nil
}

func (h *echoController) PostSpec(c echo.Context) error {
	var dto procdec.DecSpec
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	validateErr := dto.Validate()
	if validateErr != nil {
		h.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
	}
	spec, convErr := MsgToDecSpec(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	snap, createErr := h.api.Create(spec)
	if createErr != nil {
		return createErr
	}
	return c.JSON(http.StatusCreated, MsgFromDecSnap(snap))
}

func (h *echoController) GetSnap(c echo.Context) error {
	var dto sdk.SemRef
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	ref, convErr := semtype.MsgToRef(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	snap, retrieveErr := h.api.RetrieveSnap(ref)
	if retrieveErr != nil {
		return retrieveErr
	}
	return c.JSON(http.StatusOK, MsgFromDecSnap(snap))
}
