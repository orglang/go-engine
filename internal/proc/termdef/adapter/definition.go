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
	var err error
	switch s := spec.(type) {
	case termexp.DescCloseSpec:
		err = out.FromProcTermDefCloseSpec(webapi.ProcTermDefCloseSpec{ContPh: symbol.ConvertToString(s.ContChnlPH)})
	case termexp.DescWaitSpec:
		var cont webapi.ProcTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil {
			err = out.FromProcTermDefWaitSpec(webapi.ProcTermDefWaitSpec{ContPh: symbol.ConvertToString(s.ContChnlPH), ContExpr: cont})
		}
	case termexp.DescSendSpec:
		err = out.FromProcTermDefSendSpec(webapi.ProcTermDefSendSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ValPh: symbol.ConvertToString(s.ValChnlPH)})
	case termexp.DescRecvSpec:
		var cont webapi.ProcTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil {
			err = out.FromProcTermDefRecvSpec(webapi.ProcTermDefRecvSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), NewPh: symbol.ConvertToString(s.NewChnlPH), ContExpr: cont})
		}
	case termexp.DescLabSpec:
		var cont webapi.ProcTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil {
			err = out.FromProcTermDefLabSpec(webapi.ProcTermDefLabSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ValLabQn: uniqsym.ConvertToString(s.ValLabQN), ContExpr: cont})
		}
	case termexp.DescCaseSpec:
		branches := make(map[string]webapi.ProcTermDefinitionSpec, len(s.ContExps))
		for label, child := range s.ContExps {
			var cont webapi.ProcTermDefinitionSpec
			cont, err = DefinitionToGenerated(child)
			if err != nil {
				break
			}
			branches[uniqsym.ConvertToString(label)] = cont
		}
		if err == nil {
			err = out.FromProcTermDefCaseSpec(webapi.ProcTermDefCaseSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExps: branches})
		}
	case termexp.DescLinkSpec:
		ys := make([]string, len(s.Ys))
		for i, y := range s.Ys {
			ys[i] = identity.ConvertToString(y)
		}
		err = out.FromProcTermDefLinkSpec(webapi.ProcTermDefLinkSpec{ProcQn: uniqsym.ConvertToString(s.ProcTermQN), X: identity.ConvertToString(s.X), Ys: ys})
	case termexp.DescFwdSpec:
		err = out.FromProcTermDefFwdSpec(webapi.ProcTermDefFwdSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContPh: symbol.ConvertToString(s.ContChnlPH)})
	case termexp.DescCallSpec:
		var cont webapi.ProcTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil {
			chs := make([]string, len(s.ValChnlPHs))
			for i, ch := range s.ValChnlPHs { chs[i] = symbol.ConvertToString(ch) }
			err = out.FromProcTermDefCallSpec(webapi.ProcTermDefCallSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), NewPh: symbol.ConvertToString(s.NewChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ValPhs: chs, ContExpr: cont})
		}
	case termexp.DescSpawnSpec:
		var cont webapi.ProcTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil {
			chs := make([]string, len(s.NewChnlPHs))
			for i, ch := range s.NewChnlPHs { chs[i] = symbol.ConvertToString(ch) }
			err = out.FromProcTermDefSpawnSpec(webapi.ProcTermDefSpawnSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), NewPhs: chs, ContExpr: cont})
		}
	case termexp.DescAcquireSpec:
		var cont webapi.ProcTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil { err = out.FromProcTermDefAcquireSpec(webapi.ProcTermDefAcquireSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExpr: cont}) }
	case termexp.DescAcceptSpec:
		var cont webapi.ProcTermDefinitionSpec
		cont, err = DefinitionToGenerated(s.ContExp)
		if err == nil { err = out.FromProcTermDefAcceptSpec(webapi.ProcTermDefAcceptSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH), ContExpr: cont}) }
	case termexp.DescDetachSpec:
		err = out.FromProcTermDefDetachSpec(webapi.ProcTermDefDetachSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH)})
	case termexp.DescReleaseSpec:
		err = out.FromProcTermDefReleaseSpec(webapi.ProcTermDefReleaseSpec{ConnPh: symbol.ConvertToString(s.ConnChnlPH)})
	default:
		return out, fmt.Errorf("unsupported process Definition spec %T", spec)
	}
	if err != nil { return out, err }
	return out, nil
}
