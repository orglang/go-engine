package adapter

import (
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdk "github.com/orglang/go-sdk/adt/termsem"
	sdktermdec "github.com/orglang/go-sdk/proc/termdec"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

// Server-side primary adapter
type controllerEcho struct {
	webapi.ServerInterface
	api termdec.API
	log *slog.Logger
}

func newControllerEcho(a termdec.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{api: a, log: l.With(name)}
}

func cfgEchoController(e *echo.Echo, h *controllerEcho) error {
	e.POST("/api/v1/procs/decs", h.PostSpec)
	e.POST("/api/v1/processes/declarations", h.CreateProcessDeclaration)
	e.GET("/api/v1/procs/decs/:id", h.GetSnap)
	return nil
}

func (h *controllerEcho) PostSpec(c echo.Context) error {
	var dto sdktermdec.DecSpec
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	validateErr := dto.Validate()
	if validateErr != nil {
		h.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
	}
	spec, convErr := MsgToDecSpec(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	snap, createErr := h.api.Create(spec)
	if createErr != nil {
		return createErr
	}
	return c.JSON(http.StatusCreated, MsgFromDecSnap(snap))
}

func (h *controllerEcho) GetSnap(c echo.Context) error {
	var dto sdk.SemRef
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
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

	spec, err := decSpecFromGenerated(dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	snap, err := h.api.Create(spec)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, webapi.AnyTermRef{
		TermId: identity.ConvertToString(snap.TermRef.TermID),
		TermRn: int64(snap.TermRef.TermRN),
	})
}

func decSpecFromGenerated(dto webapi.ProcDecSpec) (termdec.DecSpec, error) {
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

func connVarFromGenerated(dto webapi.AnyConnVarSpec) (termvar.VarSpec, error) {
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
