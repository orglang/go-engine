package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/identity"
	webapi "orglang/go-engine/internal/app/web/api"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

type Controller struct { api termdec.API; log *slog.Logger }

func NewControllerEcho(api termdec.API, l *slog.Logger) *Controller {
	name := slog.String("name", reflect.TypeFor[Controller]().Name())
	return &Controller{api: api, log: l.With(name)}
}

func (h *Controller) GetProcessDeclaration(c *echo.Context, id string, _ webapi.GetProcessDeclarationParams) error {
	termID, err := identity.ConvertFromString(id)
	if err != nil { h.log.Error("conversion failed", slog.Any("id", id), slog.Any("err", err)); return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }
	snap, err := h.api.RetrieveSnap(c.Request().Context(), termID)
	if err != nil { return err }
	return c.JSON(http.StatusOK, DecSnapToGenerated(snap))
}

func (h *Controller) CreateProcessDeclaration(c *echo.Context) error {
	var dto webapi.ProcDecSpec
	if err := c.Bind(&dto); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }
	spec, err := DecSpecFromGenerated(dto)
	if err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }
	snap, err := h.api.Create(c.Request().Context(), spec)
	if err != nil { return err }
	return c.JSON(http.StatusCreated, DecSnapToGenerated(snap))
}
