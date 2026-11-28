package adapter

import (
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	typedef "orglang/go-engine/internal/pool/typedef/core"
)

type GeneratedController interface {
	CreatePoolTypeDefinition(*echo.Context) error
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

func (h *controllerEcho) CreatePoolTypeDefinition(c *echo.Context) error {
	var dto webapi.PoolTypeDefSpec
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	qn, err := uniqsym.ConvertFromString(dto.Qn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	expr, err := poolTypeExprToCore(dto.Expr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	typeID, progRef, err := h.api.Create(c.Request().Context(), typedef.DefSpec{TypeQN: qn, TypeExp: expr})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, webapi.AnyTypeRef{
		TypeId: identity.ConvertToString(typeID),
		DescRn: seqnum.ConvertToInt(progRef.ProgRN),
	})
}

func poolTypeExprToCore(dto webapi.PoolTypeExprSpec) (typedef.ExpSpec, error) {
	kind, err := dto.Discriminator()
	if err != nil {
		return nil, fmt.Errorf("read pool type expression kind: %w", err)
	}
	switch kind {
	case "one":
		return typedef.OneSpec{}, nil
	case "link":
		spec, err := dto.AsPoolTypeLinkSpec()
		if err != nil {
			return nil, err
		}
		qn, err := uniqsym.ConvertFromString(spec.Link.TypeQn)
		if err != nil {
			return nil, err
		}
		return typedef.LinkSpec{TypeQN: qn}, nil
	case "plus":
		spec, err := dto.AsPoolTypePlusSpec()
		if err != nil {
			return nil, err
		}
		qns, err := poolTypeQNsToCore(spec.ProcQns)
		if err != nil {
			return nil, err
		}
		cont, err := poolTypeExprToCore(spec.ContExpr)
		if err != nil {
			return nil, err
		}
		return typedef.PlusSpec{ProcQNs: qns, ContExp: cont}, nil
	case "with":
		spec, err := dto.AsPoolTypeWithSpec()
		if err != nil {
			return nil, err
		}
		qns, err := poolTypeQNsToCore(spec.ProcQns)
		if err != nil {
			return nil, err
		}
		cont, err := poolTypeExprToCore(spec.ContExpr)
		if err != nil {
			return nil, err
		}
		return typedef.WithSpec{ProcQNs: qns, ContExp: cont}, nil
	case "up":
		spec, err := dto.AsPoolTypeUpSpec()
		if err != nil {
			return nil, err
		}
		cont, err := poolTypeExprToCore(spec.ContExpr)
		if err != nil {
			return nil, err
		}
		return typedef.UpSpec{ContExp: cont}, nil
	case "down":
		spec, err := dto.AsPoolTypeDownSpec()
		if err != nil {
			return nil, err
		}
		cont, err := poolTypeExprToCore(spec.ContExpr)
		if err != nil {
			return nil, err
		}
		return typedef.DownSpec{ContExp: cont}, nil
	default:
		return nil, fmt.Errorf("unsupported pool type expression kind %q", kind)
	}
}

func poolTypeQNsToCore(values []string) ([]uniqsym.ADT, error) {
	result := make([]uniqsym.ADT, len(values))
	for i, value := range values {
		qn, err := uniqsym.ConvertFromString(value)
		if err != nil {
			return nil, err
		}
		result[i] = qn
	}
	return result, nil
}
