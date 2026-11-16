package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdk1 "github.com/orglang/go-sdk/adt/compsem"
	sdk2 "github.com/orglang/go-sdk/proc/compstep"

	"orglang/go-engine/internal/lib/lf"

	"orglang/go-engine/internal/adt/compsem"
	webapi "orglang/go-engine/internal/app/web/api"
	compexec "orglang/go-engine/internal/proc/compexec/core"
	compstep "orglang/go-engine/internal/proc/compstep/adapter"
)

// Controller is the inbound HTTP port implemented by the Echo adapter.
type Controller interface {
	PostStep(echo.Context) error
}

func cfgLegacyStepRoute(e *echo.Echo, h *controllerEcho) {
	e.POST("/api/v1/procs/:id/steps", h.PostStep)
}

// Server-side primary adapter
type controllerEcho struct {
	webapi.ServerInterface
	api compexec.API
	log *slog.Logger
}

var _ webapi.ServerInterface = (*controllerEcho)(nil)

func newControllerEcho(a compexec.API, l *slog.Logger) *controllerEcho {
	return &controllerEcho{api: a, log: l}
}

func (h *controllerEcho) GetProcessExecution(c echo.Context, id string, params webapi.GetProcessExecutionParams) error {
	dto := sdk1.SemRef{CompID: id}
	if params.Rn != nil {
		dto.CompRN = *params.Rn
	}
	ref, convErr := compsem.MsgToRef(dto)
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

func (h *controllerEcho) PostStep(c echo.Context) error {
	var dto sdk2.StepSpec
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
