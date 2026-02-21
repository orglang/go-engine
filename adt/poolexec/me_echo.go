package poolexec

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdk "github.com/orglang/go-sdk/adt/implsem"
	"github.com/orglang/go-sdk/adt/poolexec"

	"orglang/go-engine/lib/te"

	"orglang/go-engine/adt/implsem"
)

// Server-side primary adapter
type echoController struct {
	api API
	ssr te.Renderer
	log *slog.Logger
}

func newEchoController(a API, r te.Renderer, l *slog.Logger) *echoController {
	name := slog.String("name", reflect.TypeFor[echoController]().Name())
	return &echoController{a, r, l.With(name)}
}

func cfgEchoController(e *echo.Echo, h *echoController) error {
	e.POST("/api/v1/pools/execs", h.PostSpec)
	e.GET("/api/v1/pools/execs/:id", h.GetSnap)
	return nil
}

func (h *echoController) PostSpec(c echo.Context) error {
	var dto poolexec.ExecSpec
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
	spec, convertErr := MsgToExecSpec(dto)
	if convertErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convertErr
	}
	ref, createErr := h.api.Run(spec)
	if createErr != nil {
		return createErr
	}
	return c.JSON(http.StatusCreated, implsem.MsgFromRef(ref))
}

func (h *echoController) GetSnap(c echo.Context) error {
	var dto sdk.SemRef
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		return bindErr
	}
	ref, convertErr := implsem.MsgToRef(dto)
	if convertErr != nil {
		return convertErr
	}
	snap, retrieveErr := h.api.RetrieveSnap(ref)
	if retrieveErr != nil {
		return retrieveErr
	}
	return c.JSON(http.StatusOK, MsgFromExecSnap(snap))
}
