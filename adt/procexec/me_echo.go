package procexec

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"github.com/orglang/go-sdk/adt/procexec"
	sdk "github.com/orglang/go-sdk/adt/procstep"
	"github.com/orglang/go-sdk/adt/uniqref"

	"orglang/go-runtime/lib/lf"

	"orglang/go-runtime/adt/identity"
	"orglang/go-runtime/adt/procstep"
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
	e.POST("/api/v1/procs/:id/execs", h.PostExec)
	e.POST("/api/v1/pools/:id/steps", h.PostStep)
	return nil
}

func (h *echoController) GetSnap(c echo.Context) error {
	var dto uniqref.Msg
	bindingErr := c.Bind(&dto)
	if bindingErr != nil {
		h.log.Error("binding failed", slog.Any("dto", dto))
		return bindingErr
	}
	id, conversionErr := identity.ConvertFromString(dto.ID)
	if conversionErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return conversionErr
	}
	snap, retrievalErr := h.api.Retrieve(id)
	if retrievalErr != nil {
		return retrievalErr
	}
	return c.JSON(http.StatusOK, MsgFromExecSnap(snap))
}

func (h *echoController) PostExec(c echo.Context) error {
	var dto procexec.ExecSpec
	bindingErr := c.Bind(&dto)
	if bindingErr != nil {
		h.log.Error("binding failed", slog.Any("dto", dto))
		return bindingErr
	}
	spec, conversionErr := MsgToExecSpec(dto)
	if conversionErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return conversionErr
	}
	runningErr := h.api.Run(spec)
	if runningErr != nil {
		return runningErr
	}
	return c.NoContent(http.StatusOK)
}

func (h *echoController) PostStep(c echo.Context) error {
	var dto sdk.StepSpec
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
	spec, conversionErr := procstep.MsgToStepSpec(dto)
	if conversionErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return conversionErr
	}
	takingErr := h.api.Take(spec)
	if takingErr != nil {
		return takingErr
	}
	return c.NoContent(http.StatusOK)
}
