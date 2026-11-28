package adapter

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	"orglang/go-engine/internal/pool/termdec/core"
)

type GeneratedController interface {
	CreatePoolDeclaration(*echo.Context) error
}

type controllerEcho struct {
	api core.API
}

var _ GeneratedController = (*controllerEcho)(nil)

func newControllerEcho(api core.API) *controllerEcho {
	return &controllerEcho{api: api}
}

func (h *controllerEcho) CreatePoolDeclaration(c *echo.Context) error {
	var dto webapi.PoolDecSpec
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	qn, err := uniqsym.ConvertFromString(dto.DescQn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	liabVar, err := commVarFromGenerated(dto.LiabVar)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	assetVars := make([]termvar.VarSpec, len(dto.AssetVars))
	for i, variable := range dto.AssetVars {
		assetVars[i], err = commVarFromGenerated(variable)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	termID, progRef, err := h.api.Create(c.Request().Context(), core.DecSpec{TermQN: qn, LiabVar: liabVar, AssetVars: assetVars})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, webapi.AnyTermRef{
		TermId: identity.ConvertToString(termID),
		DescRn: seqnum.ConvertToInt(progRef.ProgRN),
	})
}

func commVarFromGenerated(dto webapi.AnyCommVarSpec) (termvar.VarSpec, error) {
	channel, err := symbol.ConvertFromString(dto.CommPh)
	if err != nil {
		return termvar.VarSpec{}, err
	}
	typeQN, err := uniqsym.ConvertFromString(dto.TypeQn)
	if err != nil {
		return termvar.VarSpec{}, err
	}
	return termvar.VarSpec{ChnlPH: channel, TypeQN: typeQN}, nil
}
