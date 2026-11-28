package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/app/web/api"
	execeval "orglang/go-engine/internal/pool/execeval/core"
	execstep "orglang/go-engine/internal/pool/execstep/core"
	execvar "orglang/go-engine/internal/pool/execvar/core"
)

// Controller is the inbound HTTP port implemented by the Echo adapter.
type Controller interface {
	PostSpec(*echo.Context) error
	PostSpec2(*echo.Context) error
	PostSpec3(*echo.Context) error
}

// Server-side primary adapter.
type controllerEcho struct {
	api execeval.API
	log *slog.Logger
}

func newControllerEcho(evalAPI execeval.API, log *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{evalAPI, log.With(name)}
}

// RegisterRoutes attaches the pool execution HTTP endpoints to an
// Echo server. Controller construction is owned by the application container.
func RegisterRoutes(server *echo.Echo, controller Controller) error {
	server.POST("/api/v1/pools/executions", controller.PostSpec)
	server.POST("/api/v1/pools/executions/steps", controller.PostSpec2)
	server.POST("/api/v1/pools/executions/spawns", controller.PostSpec3)
	return nil
}

func (c *controllerEcho) PostSpec(ctx *echo.Context) error {
	var dto api.AnyEvalSpec
	if err := ctx.Bind(&dto); err != nil {
		c.log.Error("binding failed", slog.Any("dto", reflect.TypeFor[api.AnyEvalSpec]()))
		return err
	}
	spec, err := poolExecSpecToCore(dto)
	if err != nil {
		c.log.Error("conversion failed", slog.Any("dto", dto), slog.Any("error", err))
		return err
	}
	ref, err := c.api.Run(ctx.Request().Context(), spec)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, rootRefToAPI(ref))
}

func poolExecSpecToCore(dto api.AnyEvalSpec) (execeval.EvalSpec, error) {
	termQN, err := uniqsym.ConvertFromString(dto.TermQn)
	if err != nil {
		return execeval.EvalSpec{}, err
	}
	liabVar, err := poolExecVarToCore(dto.LiabVar)
	if err != nil {
		return execeval.EvalSpec{}, err
	}
	assetVars := make([]execvar.VarSpec, len(dto.AssetVars))
	for i, assetVar := range dto.AssetVars {
		assetVars[i], err = poolExecVarToCore(assetVar)
		if err != nil {
			return execeval.EvalSpec{}, fmt.Errorf("asset_vars[%d]: %w", i, err)
		}
	}
	return execeval.EvalSpec{TermQN: termQN, LiabVar: liabVar, AssetVars: assetVars}, nil
}

func poolExecVarToCore(dto api.AnyEvalVarSpec) (execvar.VarSpec, error) {
	if dto.CommPh == "" {
		return execvar.VarSpec{}, errors.New("comm_ph cannot be blank")
	}
	commPH, err := symbol.ConvertFromString(dto.CommPh)
	if err != nil {
		return execvar.VarSpec{}, err
	}
	termQN, err := uniqsym.ConvertFromString(dto.TermQn)
	if err != nil {
		return execvar.VarSpec{}, err
	}
	return execvar.VarSpec{ChnlPH: commPH, TermQN: termQN}, nil
}

func (c *controllerEcho) PostSpec2(ctx *echo.Context) error {
	var dto api.PoolStepSpec
	if err := ctx.Bind(&dto); err != nil {
		c.log.Error("binding failed", slog.Any("dto", reflect.TypeFor[api.PoolStepSpec]()))
		return err
	}
	spec, err := poolStepSpecToCore(dto)
	if err != nil {
		c.log.Error("conversion failed", slog.Any("dto", dto), slog.Any("error", err))
		return err
	}
	if err = c.api.Take(ctx.Request().Context(), spec); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (c *controllerEcho) PostSpec3(ctx *echo.Context) error {
	var dto api.PoolStepSpec
	if err := ctx.Bind(&dto); err != nil {
		c.log.Error("binding failed", slog.Any("dto", reflect.TypeFor[api.PoolStepSpec]()))
		return err
	}
	spec, err := poolStepSpecToCore(dto)
	if err != nil {
		c.log.Error("conversion failed", slog.Any("dto", dto), slog.Any("error", err))
		return err
	}
	ref, err := c.api.Spawn(ctx.Request().Context(), spec)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, rootRefToAPI(ref))
}

func poolStepSpecToCore(dto api.PoolStepSpec) (execstep.StepSpec, error) {
	ref, err := rootRefFromAPI(dto.RootRef)
	if err != nil {
		return execstep.StepSpec{}, err
	}
	expr, err := poolTermExprToCore(dto.PoolExprr)
	if err != nil {
		return execstep.StepSpec{}, err
	}
	return execstep.StepSpec{RootRef: ref, PoolExpr: exp}, nil
}

func rootRefFromAPI(dto api.AnyRootRef) (execroot.RootRef, error) {
	id, err := identity.ConvertFromString(dto.ExecId)
	if err != nil {
		return execroot.RootRef{}, err
	}
	return execroot.RootRef{ExecID: id, ExecRN: seqnum.ConvertFromInt(dto.ExecRn)}, nil
}

func rootRefToAPI(ref execroot.RootRef) api.AnyRootRef {
	return api.AnyRootRef{
		ExecId: identity.ConvertToString(ref.ExecID),
		ExecRn: seqnum.ConvertToInt(ref.ExecRN),
	}
}

func termRefFromAPI(dto api.AnyTermRef) (identity.ADT, seqnum.ADT, error) {
	id, err := identity.ConvertFromString(dto.TermId)
	if err != nil {
		return identity.ADT{}, 0, err
	}
	return id, seqnum.ConvertFromInt(dto.ProgRn), nil
}

func poolTermExprToCore(dto api.PoolTermExprSpec) (execeval.ExpSpec, error) {
	data, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	var discriminator struct {
		Kind string `json:"kind"`
	}
	if err = json.Unmarshal(data, &discriminator); err != nil {
		return nil, err
	}
	if discriminator.Kind == "" {
		if string(data) == "null" || string(data) == "{}" {
			return nil, nil
		}
		return nil, errors.New("pool_expr: kind cannot be blank")
	}

	switch discriminator.Kind {
	case "acquire":
		value, err := dto.AsPoolTermAcquireSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.CommPh)
		if err != nil {
			return nil, err
		}
		cont, err := poolTermExprToCore(value.ContExpr)
		if err != nil {
			return nil, err
		}
		return execeval.AcquireSpec{CommChnlPH: ph, ContExp: cont}, nil
	case "accept":
		value, err := dto.AsPoolTermAcceptSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.CommPh)
		if err != nil {
			return nil, err
		}
		cont, err := poolTermExprToCore(value.ContExpr)
		if err != nil {
			return nil, err
		}
		return execeval.AcceptSpec{CommChnlPH: ph, ContExp: cont}, nil
	case "hire":
		value, err := dto.AsPoolTermHireSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.CommPh)
		if err != nil {
			return nil, err
		}
		qn, err := uniqsym.ConvertFromString(value.ProcQn)
		if err != nil {
			return nil, err
		}
		cont, err := poolTermExprToCore(value.ContExpr)
		if err != nil {
			return nil, err
		}
		return execeval.HireSpec{CommChnlPH: ph, ProcTermQN: qn, ContExp: cont}, nil
	case "apply":
		value, err := dto.AsPoolTermApplySpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.CommPh)
		if err != nil {
			return nil, err
		}
		qn, err := uniqsym.ConvertFromString(value.ProcQn)
		if err != nil {
			return nil, err
		}
		cont, err := poolTermExprToCore(value.ContExpr)
		if err != nil {
			return nil, err
		}
		return execeval.ApplySpec{CommChnlPH: ph, ProcTermQN: qn, ContExp: cont}, nil
	case "release":
		value, err := dto.AsPoolTermReleaseSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.CommPh)
		if err != nil {
			return nil, err
		}
		return execeval.ReleaseSpec{CommChnlPH: ph}, nil
	case "detach":
		value, err := dto.AsPoolTermDetachSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.CommPh)
		if err != nil {
			return nil, err
		}
		return execeval.DetachSpec{CommChnlPH: ph}, nil
	case "fire":
		value, err := dto.AsPoolTermFireSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.CommPh)
		if err != nil {
			return nil, err
		}
		return execeval.FireSpec{CommChnlPH: ph}, nil
	case "quit":
		value, err := dto.AsPoolTermQuitSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.CommPh)
		if err != nil {
			return nil, err
		}
		qn, err := uniqsym.ConvertFromString(value.ProcQn)
		if err != nil {
			return nil, err
		}
		return execeval.QuitSpec{CommChnlPH: ph, ProcTermQN: qn}, nil
	case "spawn":
		value, err := dto.AsPoolTermSpawnSpec()
		if err != nil {
			return nil, err
		}
		termID, termRN, err := termRefFromAPI(value.ProcDr)
		if err != nil {
			return nil, err
		}
		procExecRefs := make([]execroot.RootRef, len(value.ProcExecRefs))
		for i, item := range value.ProcExecRefs {
			procExecRefs[i], err = rootRefFromAPI(item)
			if err != nil {
				return nil, fmt.Errorf("proc_exec_refs[%d]: %w", i, err)
			}
		}
		return execeval.SpawnSpec2{ProcTermID: termID, ProcTermRN: termRN, ProcExecRefs: procExecRefs}, nil
	default:
		return nil, fmt.Errorf("pool_expr: unsupported kind %q", discriminator.Kind)
	}
}
