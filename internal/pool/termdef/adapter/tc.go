package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

// DefinitionToGenerated converts a persisted pool Definition to its public TypeSpec DTO.
func DefinitionToGenerated(spec termexp.DescSpec) (webapi.PoolTermDefinitionSpec, error) {
	var out webapi.PoolTermDefinitionSpec

	switch s := spec.(type) {
	case termexp.DescAcquireSpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		return out, out.FromPoolTermDefAcquireSpec(webapi.PoolTermDefAcquireSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExpr: cont})
	case termexp.DescAcceptSpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		return out, out.FromPoolTermDefAcceptSpec(webapi.PoolTermDefAcceptSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExpr: cont})
	case termexp.DescHireSpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		return out, out.FromPoolTermDefHireSpec(webapi.PoolTermDefHireSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ContExpr: cont})
	case termexp.DescApplySpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		return out, out.FromPoolTermDefApplySpec(webapi.PoolTermDefApplySpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ContExpr: cont})
	case termexp.DescReleaseSpec:
		return out, out.FromPoolTermDefReleaseSpec(webapi.PoolTermDefReleaseSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH)})
	case termexp.DescDetachSpec:
		return out, out.FromPoolTermDefDetachSpec(webapi.PoolTermDefDetachSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH)})
	case termexp.DescFireSpec:
		return out, out.FromPoolTermDefFireSpec(webapi.PoolTermDefFireSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN)})
	case termexp.DescQuitSpec:
		return out, out.FromPoolTermDefQuitSpec(webapi.PoolTermDefQuitSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN)})
	case termexp.DescSpawnSpec:
		qns := make([]string, len(s.ProcCompQNs))
		for i, qn := range s.ProcCompQNs { qns[i] = uniqsym.ConvertToString(qn) }
		return out, out.FromPoolTermDefSpawnSpec(webapi.PoolTermDefSpawnSpec{ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ProcCompQns: qns})
	default:
		return out, fmt.Errorf("unsupported pool Definition spec %T", spec)
	}
}
