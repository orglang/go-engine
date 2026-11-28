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
	return exprToGenerated(def.Expr)
}

func exprToGenerated(spec termdef.ExprSpec) (webapi.ProcTermDefinitionSpec, error) {
	var out webapi.ProcTermDefinitionSpec

	switch s := spec.(type) {
	case termdef.ExprCloseSpec:
		return out, out.FromProcTermDefCloseSpec(webapi.ProcTermDefCloseSpec{ContPh: symbol.ConvertToString(s.ContChnlPH)})
	case termdef.ExprWaitSpec:
		cont, err := exprToGenerated(s.ContExpr)
		if err != nil {
			return out, err
		}
		return out, out.FromProcTermDefWaitSpec(webapi.ProcTermDefWaitSpec{ContPh: symbol.ConvertToString(s.ContChnlPH), ContExpr: cont})
	case termdef.ExprSendSpec:
		return out, out.FromProcTermDefSendSpec(webapi.ProcTermDefSendSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ValPh: symbol.ConvertToString(s.ValChnlPH)})
	case termdef.ExprRecvSpec:
		cont, err := exprToGenerated(s.ContExpr)
		if err != nil {
			return out, err
		}
		return out, out.FromProcTermDefRecvSpec(webapi.ProcTermDefRecvSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), NewPh: symbol.ConvertToString(s.NewChnlPH), ContExpr: cont})
	case termdef.ExprLabSpec:
		cont, err := exprToGenerated(s.ContExpr)
		if err != nil {
			return out, err
		}
		return out, out.FromProcTermDefLabSpec(webapi.ProcTermDefLabSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ValLabQn: uniqsym.ConvertToString(s.ValLabQN), ContExpr: cont})
	case termdef.ExprCaseSpec:
		branches := make(map[string]webapi.ProcTermDefinitionSpec, len(s.ContExprs))
		for label, child := range s.ContExprs {
			cont, err := exprToGenerated(child)
			if err != nil {
				return out, err
			}
			branches[uniqsym.ConvertToString(label)] = cont
		}
		return out, out.FromProcTermDefCaseSpec(webapi.ProcTermDefCaseSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ContExps: branches})
	case termdef.ExprLinkSpec:
		ys := make([]string, len(s.Ys))
		for i, y := range s.Ys {
			ys[i] = identity.ConvertToString(y)
		}
		return out, out.FromProcTermDefLinkSpec(webapi.ProcTermDefLinkSpec{ProcQn: uniqsym.ConvertToString(s.ProcTermQN), X: identity.ConvertToString(s.X), Ys: ys})
	case termdef.ExprFwdSpec:
		return out, out.FromProcTermDefFwdSpec(webapi.ProcTermDefFwdSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ContPh: symbol.ConvertToString(s.ContChnlPH)})
	case termdef.ExprCallSpec:
		cont, err := exprToGenerated(s.ContExpr)
		if err != nil {
			return out, err
		}
		chs := make([]string, len(s.ValChnlPHs))
		for i, ch := range s.ValChnlPHs {
			chs[i] = symbol.ConvertToString(ch)
		}
		return out, out.FromProcTermDefCallSpec(webapi.ProcTermDefCallSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), NewPh: symbol.ConvertToString(s.NewChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), ValPhs: chs, ContExpr: cont})
	case termdef.ExprSpawnSpec:
		cont, err := exprToGenerated(s.ContExpr)
		if err != nil {
			return out, err
		}
		chs := make([]string, len(s.NewChnlPHs))
		for i, ch := range s.NewChnlPHs {
			chs[i] = symbol.ConvertToString(ch)
		}
		return out, out.FromProcTermDefSpawnSpec(webapi.ProcTermDefSpawnSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ProcQn: uniqsym.ConvertToString(s.ProcTermQN), NewPhs: chs, ContExpr: cont})
	case termdef.ExprAcquireSpec:
		cont, err := exprToGenerated(s.ContExpr)
		if err != nil {
			return out, err
		}
		return out, out.FromProcTermDefAcquireSpec(webapi.ProcTermDefAcquireSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ContExpr: cont})
	case termdef.ExprAcceptSpec:
		cont, err := exprToGenerated(s.ContExpr)
		if err != nil {
			return out, err
		}
		return out, out.FromProcTermDefAcceptSpec(webapi.ProcTermDefAcceptSpec{CommPh: symbol.ConvertToString(s.CommChnlPH), ContExpr: cont})
	case termdef.ExprDetachSpec:
		return out, out.FromProcTermDefDetachSpec(webapi.ProcTermDefDetachSpec{CommPh: symbol.ConvertToString(s.CommChnlPH)})
	case termdef.ExprReleaseSpec:
		return out, out.FromProcTermDefReleaseSpec(webapi.ProcTermDefReleaseSpec{CommPh: symbol.ConvertToString(s.CommChnlPH)})
	default:
		return out, fmt.Errorf("unsupported process Definition spec %T", spec)
	}
}
