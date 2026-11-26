package adapter

import (
	"errors"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termsem"
	webapi "orglang/go-engine/internal/app/web/api"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	termdef "orglang/go-engine/internal/proc/termdef/core"
)

// controllerEcho is the generated HTTP adapter for process declarations.
type controllerEcho struct {
	webapi.ServerInterface
	decAPI termdec.API
	defAPI termdef.API
	log *slog.Logger
}

var _ webapi.ServerInterface = (*controllerEcho)(nil)

func newControllerEcho(decAPI termdec.API, defAPI termdef.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{decAPI: decAPI, defAPI: defAPI, log: l.With(name)}
}

func (h *controllerEcho) GetProcessDeclaration(c *echo.Context, id string, params webapi.GetProcessDeclarationParams) error {
	termID, err := identity.ConvertFromString(id)
	if err != nil {
		h.log.Error("conversion failed", slog.Any("id", id), slog.Any("err", err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ref := termsem.SemRef{TermID: termID}
	if params.Rn != nil {
		ref.DescRN = seqnum.ConvertFromInt(*params.Rn)
	}
	snap, err := h.decAPI.RetrieveSnap(c.Request().Context(), ref)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, DecSnapToGenerated(snap))
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

func (h *controllerEcho) CreateProcessDeclaration(c *echo.Context) error {
	var dto webapi.ProcDecSpec
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	spec, err := DecSpecFromGenerated(dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	snap, err := h.decAPI.Create(c.Request().Context(), spec)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, DecSnapToGenerated(snap))
}
