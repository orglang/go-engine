package adapter

import (
	"errors"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	webapi "orglang/go-engine/internal/app/web/api"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

// GeneratedController is the inbound HTTP port for pool term definitions.
type GeneratedController interface {
	CreatePoolDeclaration(*echo.Context) error
	GetPoolDefinition(*echo.Context, string, webapi.GetPoolDefinitionParams) error
}

type controllerEcho struct {
	api termdef.API
	log *slog.Logger
}

var _ GeneratedController = (*controllerEcho)(nil)

func newControllerEcho(a termdef.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{api: a, log: l.With(name)}
}

func (h *controllerEcho) GetPoolDefinition(c *echo.Context, id string, _ webapi.GetPoolDefinitionParams) error {
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

func (h *controllerEcho) CreatePoolDeclaration(c *echo.Context) error {
	var dto webapi.PoolDecSpec
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	spec, err := DecSpecFromGenerated(dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ref, err := h.api.Create(c.Request().Context(), spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, webapi.AnyTermRef{
		TermId: identity.ConvertToString(ref.TermID),
		DescRn: seqnum.ConvertToInt(ref.DescRN),
	})
}
