package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/app/web/api"
	compexec "orglang/go-engine/internal/pool/compexec/core"
	compstep "orglang/go-engine/internal/pool/compstep/core"
	compvar "orglang/go-engine/internal/pool/compvar/core"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

// Controller is the inbound HTTP port implemented by the Echo adapter.
type Controller interface {
	PostSpec(echo.Context) error
	PostSpec2(echo.Context) error
	PostSpec3(echo.Context) error
}

// Server-side primary adapter.
type controllerEcho struct {
	api compexec.API
	log *slog.Logger
}

func newControllerEcho(execAPI compexec.API, log *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{execAPI, log.With(name)}
}

// RegisterRoutes attaches the pool computation-execution HTTP endpoints to an
// Echo server. Controller construction is owned by the application container.
func RegisterRoutes(server *echo.Echo, controller Controller) error {
	server.POST("/api/v1/pools/executions", controller.PostSpec)
	server.POST("/api/v1/pools/executions/steps", controller.PostSpec2)
	server.POST("/api/v1/pools/executions/spawns", controller.PostSpec3)
	return nil
}

func (c *controllerEcho) PostSpec(ctx echo.Context) error {
	var dto api.AnyExecSpec
	if err := ctx.Bind(&dto); err != nil {
		c.log.Error("binding failed", slog.Any("dto", reflect.TypeFor[api.AnyExecSpec]()))
		return err
	}
	spec, err := poolExecSpecToCore(dto)
	if err != nil {
		c.log.Error("conversion failed", slog.Any("dto", dto), slog.Any("error", err))
		return err
	}
	ref, err := c.api.Run(spec)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, compRefToAPI(ref))
}

func poolExecSpecToCore(dto api.AnyExecSpec) (compexec.ExecSpec, error) {
	termQN, err := uniqsym.ConvertFromString(dto.TermQn)
	if err != nil {
		return compexec.ExecSpec{}, err
	}
	liabVar, err := poolExecVarToCore(dto.LiabVar)
	if err != nil {
		return compexec.ExecSpec{}, err
	}
	assetVars := make([]compvar.VarSpec, len(dto.AssetVars))
	for i, assetVar := range dto.AssetVars {
		assetVars[i], err = poolExecVarToCore(assetVar)
		if err != nil {
			return compexec.ExecSpec{}, fmt.Errorf("asset_vars[%d]: %w", i, err)
		}
	}
	return compexec.ExecSpec{TermQN: termQN, LiabVar: liabVar, AssetVars: assetVars}, nil
}

func poolExecVarToCore(dto api.AnyExecVarSpec) (compvar.VarSpec, error) {
	if dto.ConnPh == "" {
		return compvar.VarSpec{}, errors.New("conn_ph cannot be blank")
	}
	connPH, err := symbol.ConvertFromString(dto.ConnPh)
	if err != nil {
		return compvar.VarSpec{}, err
	}
	termQN, err := uniqsym.ConvertFromString(dto.TermQn)
	if err != nil {
		return compvar.VarSpec{}, err
	}
	return compvar.VarSpec{ChnlPH: connPH, TermQN: termQN}, nil
}

func (c *controllerEcho) PostSpec2(ctx echo.Context) error {
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
	if err = c.api.Take(spec); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (c *controllerEcho) PostSpec3(ctx echo.Context) error {
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
	ref, err := c.api.Spawn(spec)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, compRefToAPI(ref))
}

func poolStepSpecToCore(dto api.PoolStepSpec) (compstep.StepSpec, error) {
	ref, err := compRefFromAPI(dto.CompRef)
	if err != nil {
		return compstep.StepSpec{}, err
	}
	exp, err := poolTermExprToCore(dto.PoolExpr)
	if err != nil {
		return compstep.StepSpec{}, err
	}
	return compstep.StepSpec{CompRef: ref, PoolExp: exp}, nil
}

func compRefFromAPI(dto api.AnyCompRef) (compsem.SemRef, error) {
	id, err := identity.ConvertFromString(dto.CompId)
	if err != nil {
		return compsem.SemRef{}, err
	}
	return compsem.SemRef{CompID: id, CompRN: seqnum.ConvertFromInt(dto.CompRn)}, nil
}

func compRefToAPI(ref compsem.SemRef) api.AnyCompRef {
	return api.AnyCompRef{
		CompId: identity.ConvertToString(ref.CompID),
		CompRn: seqnum.ConvertToInt(ref.CompRN),
	}
}

func termRefFromAPI(dto api.AnyTermRef) (termsem.SemRef, error) {
	id, err := identity.ConvertFromString(dto.TermId)
	if err != nil {
		return termsem.SemRef{}, err
	}
	return termsem.SemRef{TermID: id, TermRN: seqnum.ConvertFromInt(dto.TermRn)}, nil
}

func poolTermExprToCore(dto api.PoolTermExprSpec) (termexp.ExpSpec, error) {
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
		ph, err := symbol.ConvertFromString(value.ConnPh)
		if err != nil {
			return nil, err
		}
		cont, err := poolTermExprToCore(value.ContExpr)
		if err != nil {
			return nil, err
		}
		return termexp.AcquireSpec{ConnChnlPH: ph, ContExp: cont}, nil
	case "accept":
		value, err := dto.AsPoolTermAcceptSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.ConnPh)
		if err != nil {
			return nil, err
		}
		cont, err := poolTermExprToCore(value.ContExpr)
		if err != nil {
			return nil, err
		}
		return termexp.AcceptSpec{ConnChnlPH: ph, ContExp: cont}, nil
	case "hire":
		value, err := dto.AsPoolTermHireSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.ConnPh)
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
		return termexp.HireSpec{ConnChnlPH: ph, ProcTermQN: qn, ContExp: cont}, nil
	case "apply":
		value, err := dto.AsPoolTermApplySpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.ConnPh)
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
		return termexp.ApplySpec{ConnChnlPH: ph, ProcTermQN: qn, ContExp: cont}, nil
	case "release":
		value, err := dto.AsPoolTermReleaseSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.ConnPh)
		if err != nil {
			return nil, err
		}
		return termexp.ReleaseSpec{ConnChnlPH: ph}, nil
	case "detach":
		value, err := dto.AsPoolTermDetachSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.ConnPh)
		if err != nil {
			return nil, err
		}
		return termexp.DetachSpec{ConnChnlPH: ph}, nil
	case "fire":
		value, err := dto.AsPoolTermFireSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.ConnPh)
		if err != nil {
			return nil, err
		}
		return termexp.FireSpec{ConnChnlPH: ph}, nil
	case "quit":
		value, err := dto.AsPoolTermQuitSpec()
		if err != nil {
			return nil, err
		}
		ph, err := symbol.ConvertFromString(value.ConnPh)
		if err != nil {
			return nil, err
		}
		qn, err := uniqsym.ConvertFromString(value.ProcQn)
		if err != nil {
			return nil, err
		}
		return termexp.QuitSpec{ConnChnlPH: ph, ProcTermQN: qn}, nil
	case "spawn":
		value, err := dto.AsPoolTermSpawnSpec()
		if err != nil {
			return nil, err
		}
		termRef, err := termRefFromAPI(value.ProcDr)
		if err != nil {
			return nil, err
		}
		compRefs := make([]compsem.SemRef, len(value.ProcIrs))
		for i, item := range value.ProcIrs {
			compRefs[i], err = compRefFromAPI(item)
			if err != nil {
				return nil, fmt.Errorf("proc_irs[%d]: %w", i, err)
			}
		}
		return termexp.SpawnSpec2{ProcTermRef: termRef, ProcCompRefs: compRefs}, nil
	default:
		return nil, fmt.Errorf("pool_expr: unsupported kind %q", discriminator.Kind)
	}
}
