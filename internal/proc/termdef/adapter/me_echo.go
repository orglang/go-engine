package adapter

import (
	"errors"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/identity"
	webapi "orglang/go-engine/internal/app/web/api"
	termdecadapter "orglang/go-engine/internal/proc/termdec/adapter"
	termdef "orglang/go-engine/internal/proc/termdef/core"
)

// controllerEcho composes the generated Declaration and Definition HTTP endpoints.
type controllerEcho struct {
	webapi.ServerInterface
	decController *termdecadapter.Controller
	defAPI termdef.API
	log *slog.Logger
}

var _ webapi.ServerInterface = (*controllerEcho)(nil)

func newControllerEcho(decController *termdecadapter.Controller, defAPI termdef.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{decController: decController, defAPI: defAPI, log: l.With(name)}
}

func (h *controllerEcho) GetProcessDeclaration(c *echo.Context, id string, params webapi.GetProcessDeclarationParams) error {
	return h.decController.GetProcessDeclaration(c, id, params)
}

func (h *controllerEcho) CreateProcessDeclaration(c *echo.Context) error {
	return h.decController.CreateProcessDeclaration(c)
}

func (h *controllerEcho) GetProcessDefinition(c *echo.Context, id string, _ webapi.GetProcessDefinitionParams) error {
	termID, err := identity.ConvertFromString(id)
	if err != nil {
		h.log.Error("conversion failed", slog.Any("id", id), slog.Any("err", err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	spec, err := h.defAPI.RetrieveDefinition(c.Request().Context(), termID)
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
