package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termsem"
	webapi "orglang/go-engine/internal/app/web/api"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

// Controller handles the generated HTTP endpoints for process Declarations.
type Controller struct {
	api termdec.API
	log *slog.Logger
}

func NewControllerEcho(api termdec.API, l *slog.Logger) *Controller {
	name := slog.String("name", reflect.TypeFor[Controller]().Name())
	return &Controller{api: api, log: l.With(name)}
}

func (h *Controller) GetProcessDeclaration(c *echo.Context, id string, params webapi.GetProcessDeclarationParams) error {
	termID, err := identity.ConvertFromString(id)
	if err != nil {
		h.log.Error("conversion failed", slog.Any("id", id), slog.Any("err", err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ref := termsem.SemRef{TermID: termID}
	if params.Rn != nil {
		ref.DescRN = seqnum.ConvertFromInt(*params.Rn)
	}
	snap, err := h.api.RetrieveSnap(c.Request().Context(), ref)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, DecSnapToGenerated(snap))
}

func (h *Controller) CreateProcessDeclaration(c *echo.Context) error {
	var dto webapi.ProcDecSpec
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	spec, err := DecSpecFromGenerated(dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	snap, err := h.api.Create(c.Request().Context(), spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, DecSnapToGenerated(snap))
}
