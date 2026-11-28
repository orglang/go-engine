package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

// DefinitionToGenerated converts a persisted pool Definition to its public TypeSpec DTO.
func DefinitionToGenerated(def termdef.DefSpec) (webapi.PoolTermDefinitionSpec, error) {
	return exprToGenerated(def.Expr)
}

func exprToGenerated(spec termdef.ExprSpec) (webapi.PoolTermDefinitionSpec, error) {
	var out webapi.PoolTermDefinitionSpec

	switch s := spec.(type) {
	case termdef.ExprAcquireSpec:
		cont, err := exprToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromPoolTermDefAcquireSpec(webapi.PoolTermDefAcquireSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ContExpr: cont})
	case termdef.ExprAcceptSpec:
		cont, err := exprToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromPoolTermDefAcceptSpec(webapi.PoolTermDefAcceptSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ContExpr: cont})
	case termdef.ExprHireSpec:
		cont, err := exprToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromPoolTermDefHireSpec(webapi.PoolTermDefHireSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ContExpr: cont})
	case termdef.ExprApplySpec:
		cont, err := exprToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromPoolTermDefApplySpec(webapi.PoolTermDefApplySpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ContExpr: cont})
	case termdef.ExprReleaseSpec:
		return out, out.FromPoolTermDefReleaseSpec(webapi.PoolTermDefReleaseSpec{CommPh: symbol.ConvertToString(s.CommChnlPH)})
	case termdef.ExprDetachSpec:
		return out, out.FromPoolTermDefDetachSpec(webapi.PoolTermDefDetachSpec{CommPh: symbol.ConvertToString(s.CommChnlPH)})
	case termdef.ExprFireSpec:
		return out, out.FromPoolTermDefFireSpec(webapi.PoolTermDefFireSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN)})
	case termdef.ExprQuitSpec:
		return out, out.FromPoolTermDefQuitSpec(webapi.PoolTermDefQuitSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN)})
	case termdef.ExprSpawnSpec:
		qns := make([]string, len(s.ProcCompQNs))
		for i, qn := range s.ProcCompQNs {
			qns[i] = uniqsym.ConvertToString(qn)
		}
		return out, out.FromPoolTermDefSpawnSpec(webapi.PoolTermDefSpawnSpec{ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ProcCompQns: qns})
	default:
		return out, fmt.Errorf("unsupported pool Definition spec %T", spec)
	}
}
