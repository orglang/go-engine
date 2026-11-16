package adapter

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	sdk1 "github.com/orglang/go-sdk/adt/compsem"

	"orglang/go-engine/internal/adt/compsem"
	webapi "orglang/go-engine/internal/app/web/api"
	compexec "orglang/go-engine/internal/proc/compexec/core"
)

// GeneratedController is the generated process-execution HTTP port.
type GeneratedController interface {
	webapi.ServerInterface
}

// Server-side primary adapter
type controllerEcho struct {
	webapi.ServerInterface
	api compexec.API
	log *slog.Logger
}

var _ webapi.ServerInterface = (*controllerEcho)(nil)
var _ GeneratedController = (*controllerEcho)(nil)

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

