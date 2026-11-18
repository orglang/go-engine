package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/app/web/api"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

// Controller is the inbound HTTP port implemented by the Echo adapter.
type Controller interface {
	PostSpec(echo.Context) error
}

// RegisterRoutes attaches the pool term-definition endpoint to Echo.
func RegisterRoutes(e *echo.Echo, h Controller) error {
	e.POST("/api/v1/pools/declarations", h.PostSpec)
	return nil
}

// Server-side primary adapter
type controllerEcho struct {
	api termdef.API
	log *slog.Logger
}

func newControllerEcho(a termdef.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{a, l.With(name)}
}

func cfgEchoController(e *echo.Echo, h *controllerEcho) error {
	return RegisterRoutes(e, h)
}

func (h *controllerEcho) PostSpec(c echo.Context) error {
	var dto api.PoolDecSpec
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	spec, convErr := poolDecSpecToCore(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	ref, createErr := h.api.Create(spec)
	if createErr != nil {
		return createErr
	}
	return c.JSON(http.StatusCreated, termsem.MsgFromRef(ref))
}

func poolDecSpecToCore(dto api.PoolDecSpec) (termdef.DefSpec, error) {
	termQN, err := uniqsym.ConvertFromString(dto.DescQn)
	if err != nil {
		return termdef.DefSpec{}, err
	}
	liabVar, err := poolConnVarToCore(dto.LiabVar)
	if err != nil {
		return termdef.DefSpec{}, err
	}
	assetVars := make([]termvar.VarSpec, len(dto.AssetVars))
	for i, assetVar := range dto.AssetVars {
		assetVars[i], err = poolConnVarToCore(assetVar)
		if err != nil {
			return termdef.DefSpec{}, err
		}
	}
	return termdef.DefSpec{
		TermQN:    termQN,
		LiabVar:   liabVar,
		AssetVars: assetVars,
	}, nil
}

func poolConnVarToCore(dto api.AnyConnVarSpec) (termvar.VarSpec, error) {
	typeQN, err := uniqsym.ConvertFromString(dto.TypeQn)
	if err != nil {
		return termvar.VarSpec{}, err
	}
	return termvar.VarSpec{
		ChnlPH: dto.ConnPh,
		TypeQN: typeQN,
	}, nil
}
