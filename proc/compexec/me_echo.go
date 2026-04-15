package compexec

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdkstep "github.com/orglang/go-sdk/adt/procstep"
	sdksem "github.com/orglang/go-sdk/adt/semterm"

	"orglang/go-engine/lib/lf"
	"orglang/go-engine/proc/compstep"

	"orglang/go-engine/adt/semterm"
)

// Server-side primary adapter
type echoController struct {
	api API
	log *slog.Logger
}

func newEchoController(a API, l *slog.Logger) *echoController {
	return &echoController{a, l}
}

func cfgEchoController(e *echo.Echo, h *echoController) error {
	e.GET("/api/v1/procs/:id", h.GetSnap)
	e.POST("/api/v1/procs/:id/steps", h.PostStep)
	return nil
}

func (h *echoController) GetSnap(c echo.Context) error {
	var dto sdksem.TermRef
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", dto))
		return bindErr
	}
	ref, convErr := semterm.MsgToRef(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	snap, retrieveErr := h.api.RetrieveSnap(ref)
	if retrieveErr != nil {
		return retrieveErr
	}
	return c.JSON(http.StatusOK, MsgFromExecSnap(snap))
}

func (h *echoController) PostStep(c echo.Context) error {
	var dto sdkstep.StepSpec
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
	spec, convErr := compstep.MsgToStepSpec(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	takingErr := h.api.Take(spec)
	if takingErr != nil {
		return takingErr
	}
	return c.NoContent(http.StatusOK)
}
