package procexec

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdksem "github.com/orglang/go-sdk/adt/implsem"
	sdkstep "github.com/orglang/go-sdk/adt/procstep"

	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/procstep"
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
	var dto sdksem.SemRef
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", dto))
		return bindErr
	}
	ref, convertErr := implsem.MsgToRef(dto)
	if convertErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convertErr
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
	spec, convertErr := procstep.MsgToStepSpec(dto)
	if convertErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convertErr
	}
	takingErr := h.api.Take(spec)
	if takingErr != nil {
		return takingErr
	}
	return c.NoContent(http.StatusOK)
}
