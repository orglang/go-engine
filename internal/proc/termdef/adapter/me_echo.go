package adapter

import (
	"errors"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/identity"
	webapi "orglang/go-engine/internal/app/web/api"
	termdef "orglang/go-engine/internal/proc/termdef/core"
)

// DefinitionController handles the generated HTTP endpoint for process Definitions.
type DefinitionController struct {
	api termdef.API
	log *slog.Logger
}

func NewDefinitionController(api termdef.API, l *slog.Logger) *DefinitionController {
	name := slog.String("name", reflect.TypeFor[DefinitionController]().Name())
	return &DefinitionController{api: api, log: l.With(name)}
}

func (h *DefinitionController) GetProcessDefinition(c *echo.Context, id string, _ webapi.GetProcessDefinitionParams) error {
	termID, err := identity.ConvertFromString(id)
	if err != nil {
		h.log.Error("conversion failed", slog.Any("id", id), slog.Any("err", err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	spec, err := h.api.RetrieveDefinition(c.Request().Context(), termID)
	if err != nil {
		switch {
		case errors.Is(err, termdef.ErrDefinitionMissing), errors.Is(err, termdef.ErrTermDoesNotExist):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return err
		}
	}
	dto, err := DefinitionToGenerated(spec)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dto)
}
