package adapter

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
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

func (h *controllerEcho) GetProcessExecution(c *echo.Context, id string, params webapi.GetProcessExecutionParams) error {
	compID, err := identity.ConvertFromString(id)
	if err != nil {
		h.log.Error("conversion failed", slog.Any("id", id), slog.Any("err", err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ref := compsem.SemRef{CompID: compID}
	if params.Rn != nil {
		ref.CompRN = seqnum.ConvertFromInt(*params.Rn)
	}
	snap, err := h.api.RetrieveSnap(ref)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ExecSnapToGenerated(snap))
}
