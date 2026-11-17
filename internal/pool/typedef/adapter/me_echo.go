package adapter

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"orglang/go-engine/internal/app/web/api"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/lf"
	typedef "orglang/go-engine/internal/pool/typedef/core"
	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

// Controller is the inbound HTTP port implemented by the Echo adapter.
type Controller interface {
	PostSpec(echo.Context) error
}

// RegisterRoutes attaches the pool type-definition endpoint to Echo.
func RegisterRoutes(e *echo.Echo, h Controller) error {
	e.POST("/api/v1/pools/types", h.PostSpec)
	return nil
}

// Server-side primary adapter
type controllerEcho struct {
	api typedef.API
	log *slog.Logger
}

func newControllerEcho(a typedef.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{a, l.With(name)}
}

func cfgEchoController(e *echo.Echo, h *controllerEcho) error {
	return RegisterRoutes(e, h)
}

func (h *controllerEcho) PostSpec(c echo.Context) error {
	var dto api.PoolTypeDefSpec
	if err := c.Bind(&dto); err != nil {
		h.log.Error("binding failed", slog.Any("error", err))
		return err
	}
	ctx := c.Request().Context()
	h.log.Log(ctx, lf.LevelTrace, "posting started", slog.Any("dto", dto))
	qn, err := uniqsym.ConvertFromString(dto.Qn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	expr, err := poolTypeExprToCore(dto.Expr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ref, err := h.api.Create(typedef.DefSpec{TypeQN: qn, TypeExp: expr})
	if err != nil {
		return err
	}
	h.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("ref", ref))
	return c.JSON(http.StatusCreated, typesem.MsgFromRef(ref))
}

func poolTypeExprToCore(dto api.PoolTypeExprSpec) (typeexp.ExpSpec, error) {
	kind, err := dto.Discriminator()
	if err != nil {
		return nil, fmt.Errorf("read pool type expression kind: %w", err)
	}
	switch kind {
	case "one":
		return typeexp.OneSpec{}, nil
	case "link":
		spec, err := dto.AsPoolTypeLinkSpec()
		if err != nil { return nil, err }
		qn, err := uniqsym.ConvertFromString(spec.Link.TypeQn)
		if err != nil { return nil, err }
		return typeexp.LinkSpec{TypeQN: qn}, nil
	case "plus":
		spec, err := dto.AsPoolTypePlusSpec()
		if err != nil { return nil, err }
		qns, err := poolTypeQNsToCore(spec.ProcQns)
		if err != nil { return nil, err }
		cont, err := poolTypeExprToCore(spec.ContExpr)
		if err != nil { return nil, err }
		return typeexp.PlusSpec{ProcQNs: qns, ContExp: cont}, nil
	case "with":
		spec, err := dto.AsPoolTypeWithSpec()
		if err != nil { return nil, err }
		qns, err := poolTypeQNsToCore(spec.ProcQns)
		if err != nil { return nil, err }
		cont, err := poolTypeExprToCore(spec.ContExpr)
		if err != nil { return nil, err }
		return typeexp.WithSpec{ProcQNs: qns, ContExp: cont}, nil
	case "up":
		spec, err := dto.AsPoolTypeUpSpec()
		if err != nil { return nil, err }
		cont, err := poolTypeExprToCore(spec.ContExpr)
		if err != nil { return nil, err }
		return typeexp.UpSpec{ContExp: cont}, nil
	case "down":
		spec, err := dto.AsPoolTypeDownSpec()
		if err != nil { return nil, err }
		cont, err := poolTypeExprToCore(spec.ContExpr)
		if err != nil { return nil, err }
		return typeexp.DownSpec{ContExp: cont}, nil
	default:
		return nil, fmt.Errorf("unsupported pool type expression kind %q", kind)
	}
}

func poolTypeQNsToCore(values []string) ([]uniqsym.ADT, error) {
	result := make([]uniqsym.ADT, len(values))
	for i, value := range values {
		qn, err := uniqsym.ConvertFromString(value)
		if err != nil { return nil, err }
		result[i] = qn
	}
	return result, nil
}

// Keep the generated contract DTO JSON shape explicit at the adapter boundary.
var _ = json.Valid
