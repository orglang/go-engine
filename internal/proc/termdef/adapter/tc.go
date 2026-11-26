package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	termdef "orglang/go-engine/internal/proc/termdef/core"
)

// DefinitionToGenerated converts a persisted process Definition to its public TypeSpec DTO.
func DefinitionToGenerated(def termdef.DefSpec) (webapi.ProcTermDefinitionSpec, error) {
	return expToGenerated(def.Exp)
}

func expToGenerated(spec termdef.ExpSpec) (webapi.ProcTermDefinitionSpec, error) {
	var out webapi.ProcTermDefinitionSpec

	switch s := spec.(type) {
	case termdef.ExpCloseSpec:
		return out, out.FromProcTermDefCloseSpec(webapi.ProcTermDefCloseSpec{ContPh: symbol.ConvertToString(s.ContChnlPH)})
	case termdef.ExpWaitSpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromProcTermDefWaitSpec(webapi.ProcTermDefWaitSpec{ContPh: symbol.ConvertToString(s.ContChnlPH), ContExpr: cont})
	case termdef.ExpSendSpec:
		return out, out.FromProcTermDefSendSpec(webapi.ProcTermDefSendSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ValPh: symbol.ConvertToString(s.ValChnlPH)})
	case termdef.ExpRecvSpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromProcTermDefRecvSpec(webapi.ProcTermDefRecvSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), NewPh: symbol.ConvertToString(s.NewChnlPH), ContExpr: cont})
	case termdef.ExpLabSpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromProcTermDefLabSpec(webapi.ProcTermDefLabSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ValLabQn: uniqsym.ConvertToString(s.ValLabQN), ContExpr: cont})
	case termdef.ExpCaseSpec:
		branches := make(map[string]webapi.ProcTermDefinitionSpec, len(s.ContExps))
		for label, child := range s.ContExps {
			cont, err := expToGenerated(child)
			if err != nil {
				return out, err
			}
			branches[uniqsym.ConvertToString(label)] = cont
		}
		return out, out.FromProcTermDefCaseSpec(webapi.ProcTermDefCaseSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExps: branches})
	case termdef.ExpLinkSpec:
		ys := make([]string, len(s.Ys))
		for i, y := range s.Ys {
			ys[i] = identity.ConvertToString(y)
		}
		return out, out.FromProcTermDefLinkSpec(webapi.ProcTermDefLinkSpec{ProcQn: uniqsym.ConvertToString(s.ProcTermQN), X: identity.ConvertToString(s.X), Ys: ys})
	case termdef.ExpFwdSpec:
		return out, out.FromProcTermDefFwdSpec(webapi.ProcTermDefFwdSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContPh: symbol.ConvertToString(s.ContChnlPH)})
	case termdef.ExpCallSpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		chs := make([]string, len(s.ValChnlPHs))
		for i, ch := range s.ValChnlPHs {
			chs[i] = symbol.ConvertToString(ch)
		}
		return out, out.FromProcTermDefCallSpec(webapi.ProcTermDefCallSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), NewPh: symbol.ConvertToString(s.NewChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ValPhs: chs, ContExpr: cont})
	case termdef.ExpSpawnSpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		chs := make([]string, len(s.NewChnlPHs))
		for i, ch := range s.NewChnlPHs {
			chs[i] = symbol.ConvertToString(ch)
		}
		return out, out.FromProcTermDefSpawnSpec(webapi.ProcTermDefSpawnSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), NewPhs: chs, ContExpr: cont})
	case termdef.ExpAcquireSpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromProcTermDefAcquireSpec(webapi.ProcTermDefAcquireSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExpr: cont})
	case termdef.ExpAcceptSpec:
		cont, err := expToGenerated(s.ContExp)
		if err != nil {
			return out, err
		}
		return out, out.FromProcTermDefAcceptSpec(webapi.ProcTermDefAcceptSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExpr: cont})
	case termdef.ExpDetachSpec:
		return out, out.FromProcTermDefDetachSpec(webapi.ProcTermDefDetachSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH)})
	case termdef.ExpReleaseSpec:
		return out, out.FromProcTermDefReleaseSpec(webapi.ProcTermDefReleaseSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH)})
	default:
		return out, fmt.Errorf("unsupported process Definition spec %T", spec)
	}
}
