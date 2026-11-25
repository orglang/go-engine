package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

// DefinitionToGenerated converts a persisted process Definition to its public TypeSpec DTO.
func DefinitionToGenerated(spec termexp.DescSpec) (webapi.ProcTermDefinitionSpec, error) {
	var out webapi.ProcTermDefinitionSpec

	switch s := spec.(type) {
	case termexp.DescCloseSpec:
		return out, out.FromProcTermDefCloseSpec(webapi.ProcTermDefCloseSpec{ContPh: symbol.ConvertToString(s.ContChnlPH)})
	case termexp.DescWaitSpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		return out, out.FromProcTermDefWaitSpec(webapi.ProcTermDefWaitSpec{ContPh: symbol.ConvertToString(s.ContChnlPH), ContExpr: cont})
	case termexp.DescSendSpec:
		return out, out.FromProcTermDefSendSpec(webapi.ProcTermDefSendSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ValPh: symbol.ConvertToString(s.ValChnlPH)})
	case termexp.DescRecvSpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		return out, out.FromProcTermDefRecvSpec(webapi.ProcTermDefRecvSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), NewPh: symbol.ConvertToString(s.NewChnlPH), ContExpr: cont})
	case termexp.DescLabSpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		return out, out.FromProcTermDefLabSpec(webapi.ProcTermDefLabSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ValLabQn: uniqsym.ConvertToString(s.ValLabQN), ContExpr: cont})
	case termexp.DescCaseSpec:
		branches := make(map[string]webapi.ProcTermDefinitionSpec, len(s.ContExps))
		for label, child := range s.ContExps {
			cont, err := DefinitionToGenerated(child)
			if err != nil { return out, err }
			branches[uniqsym.ConvertToString(label)] = cont
		}
		return out, out.FromProcTermDefCaseSpec(webapi.ProcTermDefCaseSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExps: branches})
	case termexp.DescLinkSpec:
		ys := make([]string, len(s.Ys))
		for i, y := range s.Ys { ys[i] = identity.ConvertToString(y) }
		return out, out.FromProcTermDefLinkSpec(webapi.ProcTermDefLinkSpec{ProcQn: uniqsym.ConvertToString(s.ProcTermQN), X: identity.ConvertToString(s.X), Ys: ys})
	case termexp.DescFwdSpec:
		return out, out.FromProcTermDefFwdSpec(webapi.ProcTermDefFwdSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContPh: symbol.ConvertToString(s.ContChnlPH)})
	case termexp.DescCallSpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		chs := make([]string, len(s.ValChnlPHs))
		for i, ch := range s.ValChnlPHs { chs[i] = symbol.ConvertToString(ch) }
		return out, out.FromProcTermDefCallSpec(webapi.ProcTermDefCallSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), NewPh: symbol.ConvertToString(s.NewChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ValPhs: chs, ContExpr: cont})
	case termexp.DescSpawnSpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		chs := make([]string, len(s.NewChnlPHs))
		for i, ch := range s.NewChnlPHs { chs[i] = symbol.ConvertToString(ch) }
		return out, out.FromProcTermDefSpawnSpec(webapi.ProcTermDefSpawnSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), NewPhs: chs, ContExpr: cont})
	case termexp.DescAcquireSpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		return out, out.FromProcTermDefAcquireSpec(webapi.ProcTermDefAcquireSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExpr: cont})
	case termexp.DescAcceptSpec:
		cont, err := DefinitionToGenerated(s.ContExp)
		if err != nil { return out, err }
		return out, out.FromProcTermDefAcceptSpec(webapi.ProcTermDefAcceptSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExpr: cont})
	case termexp.DescDetachSpec:
		return out, out.FromProcTermDefDetachSpec(webapi.ProcTermDefDetachSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH)})
	case termexp.DescReleaseSpec:
		return out, out.FromProcTermDefReleaseSpec(webapi.ProcTermDefReleaseSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH)})
	default:
		return out, fmt.Errorf("unsupported process Definition spec %T", spec)
	}
}
