package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdk "github.com/orglang/go-sdk/adt/termsem"

	"orglang/go-engine/internal/adt/termsem"
	webapi "orglang/go-engine/internal/app/web/api"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

// controllerEcho is the generated HTTP adapter for process declarations.
type controllerEcho struct {
	webapi.ServerInterface
	api termdec.API
	log *slog.Logger
}

var _ webapi.ServerInterface = (*controllerEcho)(nil)

func newControllerEcho(a termdec.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{api: a, log: l.With(name)}
}

func (h *controllerEcho) GetProcessDeclaration(c echo.Context, id string, params webapi.GetProcessDeclarationParams) error {
	dto := sdk.SemRef{TermID: id}
	if params.Rn != nil {
		dto.TermRN = *params.Rn
	}
	return h.retrieveSnap(c, dto)
}

func (h *controllerEcho) retrieveSnap(c echo.Context, dto sdk.SemRef) error {
	ref, convErr := termsem.MsgToRef(dto)
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

func (h *controllerEcho) CreateProcessDeclaration(c echo.Context) error {
	var dto webapi.ProcDecSpec
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	spec, err := DecSpecFromGenerated(dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	snap, err := h.api.Create(spec)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, DecSnapToGenerated(snap))
}
