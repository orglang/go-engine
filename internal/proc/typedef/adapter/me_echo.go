package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/typesem"
	webapi "orglang/go-engine/internal/app/web/api"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type GeneratedController interface {
	CreateProcessTypeDefinition(echo.Context) error
	GetProcessTypeDefinition(echo.Context, string, webapi.GetProcessTypeDefinitionParams) error
}

type controllerEcho struct {
	api typedef.API
	log *slog.Logger
}

var _ GeneratedController = (*controllerEcho)(nil)

func newControllerEcho(a typedef.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{api: a, log: l.With(name)}
}

func (h *controllerEcho) CreateProcessTypeDefinition(c echo.Context) error {
	var dto webapi.ProcTypeDefSpec
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := (DefSpecVP{TypeQN: dto.TypeQn}).Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	spec, err := DefSpecFromGenerated(dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	snap, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, DefSnapToGenerated(snap).Ref)
}

func (h *controllerEcho) GetProcessTypeDefinition(c echo.Context, id string, params webapi.GetProcessTypeDefinitionParams) error {
	typeID, err := identity.ConvertFromString(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ref := typesem.SemRef{TypeID: typeID}
	if params.Rn != nil {
		ref.TypeRN = seqnum.ConvertFromInt(*params.Rn)
	}
	snap, err := h.api.RetrieveSnap(ref)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, DefSnapToGenerated(snap))
}
