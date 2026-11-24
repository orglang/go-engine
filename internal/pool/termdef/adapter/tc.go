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
	var err error
	switch s := spec.(type) {
	case termexp.DescAcquireSpec:
		var cont webapi.PoolTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil { err = out.FromPoolTermDefAcquireSpec(webapi.PoolTermDefAcquireSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExpr: cont}) }
	case termexp.DescAcceptSpec:
		var cont webapi.PoolTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil { err = out.FromPoolTermDefAcceptSpec(webapi.PoolTermDefAcceptSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExpr: cont}) }
	case termexp.DescHireSpec:
		var cont webapi.PoolTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil { err = out.FromPoolTermDefHireSpec(webapi.PoolTermDefHireSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ContExpr: cont}) }
	case termexp.DescApplySpec:
		var cont webapi.PoolTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil { err = out.FromPoolTermDefApplySpec(webapi.PoolTermDefApplySpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ContExpr: cont}) }
	case termexp.DescReleaseSpec:
		err = out.FromPoolTermDefReleaseSpec(webapi.PoolTermDefReleaseSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH)})
	case termexp.DescDetachSpec:
		err = out.FromPoolTermDefDetachSpec(webapi.PoolTermDefDetachSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH)})
	case termexp.DescFireSpec:
		err = out.FromPoolTermDefFireSpec(webapi.PoolTermDefFireSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN)})
	case termexp.DescQuitSpec:
		err = out.FromPoolTermDefQuitSpec(webapi.PoolTermDefQuitSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN)})
	case termexp.DescSpawnSpec:
		qns := make([]string, len(s.ProcCompQNs))
		for i, qn := range s.ProcCompQNs { qns[i] = uniqsym.ConvertToString(qn) }
		err = out.FromPoolTermDefSpawnSpec(webapi.PoolTermDefSpawnSpec{ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ProcCompQns: qns})
	default:
		return out, fmt.Errorf("unsupported pool Definition spec %T", spec)
	}
	if err != nil { return out, err }
	return out, nil
}
