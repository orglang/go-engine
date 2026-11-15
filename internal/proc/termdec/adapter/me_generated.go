package adapter

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/app/web/api"

	termdec "orglang/go-engine/internal/proc/termdec/core"
)

func (h *controllerEcho) CreateProcessDeclaration(c echo.Context) error {
	var dto api.ProcDecSpec
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	spec, err := decSpecFromGenerated(dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	snap, err := h.api.Create(spec)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, api.AnyTermRef{
		TermId: identity.ConvertToString(snap.TermRef.TermID),
		TermRn: int64(snap.TermRef.TermRN),
	})
}

func decSpecFromGenerated(dto api.ProcDecSpec) (termdec.DecSpec, error) {
	termQN, err := uniqsym.ConvertFromString(dto.TermQn)
	if err != nil {
		return termdec.DecSpec{}, fmt.Errorf("term_qn: %w", err)
	}

	liabVar, err := connVarFromGenerated(dto.LiabVar)
	if err != nil {
		return termdec.DecSpec{}, fmt.Errorf("liab_var: %w", err)
	}

	assetVars := make([]termvar.VarSpec, len(dto.AssetVars))
	for i, dtoVar := range dto.AssetVars {
		assetVars[i], err = connVarFromGenerated(dtoVar)
		if err != nil {
			return termdec.DecSpec{}, fmt.Errorf("asset_vars[%d]: %w", i, err)
		}
	}

	return termdec.DecSpec{
		TermQN:    termQN,
		LiabVar:   liabVar,
		AssetVars: assetVars,
	}, nil
}

func connVarFromGenerated(dto api.AnyConnVarSpec) (termvar.VarSpec, error) {
	chnlPH, err := symbol.ConvertFromString(dto.ConnPh)
	if err != nil {
		return termvar.VarSpec{}, fmt.Errorf("conn_ph: %w", err)
	}
	typeQN, err := uniqsym.ConvertFromString(dto.TypeQn)
	if err != nil {
		return termvar.VarSpec{}, fmt.Errorf("type_qn: %w", err)
	}

	return termvar.VarSpec{ChnlPH: chnlPH, TypeQN: typeQN}, nil
}
