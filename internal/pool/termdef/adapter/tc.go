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
	return expToGenerated(def.Exp)
}

func expToGenerated(spec termdef.ExpSpec) (webapi.PoolTermDefinitionSpec, error) {
	var out webapi.PoolTermDefinitionSpec

	switch s := spec.(type) {
	case termdef.ExpAcquireSpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromPoolTermDefAcquireSpec(webapi.PoolTermDefAcquireSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ContExpr: cont})
	case termdef.ExpAcceptSpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromPoolTermDefAcceptSpec(webapi.PoolTermDefAcceptSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ContExpr: cont})
	case termdef.ExpHireSpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromPoolTermDefHireSpec(webapi.PoolTermDefHireSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ContExpr: cont})
	case termdef.ExpApplySpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromPoolTermDefApplySpec(webapi.PoolTermDefApplySpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ContExpr: cont})
	case termdef.ExpReleaseSpec:
		return out, out.FromPoolTermDefReleaseSpec(webapi.PoolTermDefReleaseSpec{CommPh: symbol.ConvertToString(s.CommChnlPH)})
	case termdef.ExpDetachSpec:
		return out, out.FromPoolTermDefDetachSpec(webapi.PoolTermDefDetachSpec{CommPh: symbol.ConvertToString(s.CommChnlPH)})
	case termdef.ExpFireSpec:
		return out, out.FromPoolTermDefFireSpec(webapi.PoolTermDefFireSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN)})
	case termdef.ExpQuitSpec:
		return out, out.FromPoolTermDefQuitSpec(webapi.PoolTermDefQuitSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN)})
	case termdef.ExpSpawnSpec:
		qns := make([]string, len(s.ProcCompQNs))
		for i, qn := range s.ProcCompQNs {
			qns[i] = uniqsym.ConvertToString(qn)
		}
		return out, out.FromPoolTermDefSpawnSpec(webapi.PoolTermDefSpawnSpec{ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ProcCompQns: qns})
	default:
		return out, fmt.Errorf("unsupported pool Definition spec %T", spec)
	}
}
