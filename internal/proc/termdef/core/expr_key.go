package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
)

// ExprContentKey hashes only the declarative Definition model.
func ExprContentKey(spec ExprSpec) (contkey.ADT, []byte, error) {
	value, err := exprContentValue(spec)
	if err != nil {
		return contkey.Zero, nil, err
	}
	return contkey.HashJSON(value)
}

func exprContentValue(spec ExprSpec) (any, error) {
	child := func(s ExprSpec) (int64, error) {
		k, _, err := ExprContentKey(s)
		return int64(k), err
	}
	switch s := spec.(type) {
	case ExprCloseSpec:
		return map[string]any{"kind": "close", "cont_chnl_ph": string(s.ContChnlPH)}, nil
	case ExprWaitSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "wait", "cont_chnl_ph": string(s.ContChnlPH), "cont_expr_ck": k}, e
	case ExprSendSpec:
		return map[string]any{"kind": "send", "comm_chnl_ph": string(s.CommChnlPH), "val_chnl_ph": string(s.ValChnlPH)}, nil
	case ExprRecvSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "recv", "comm_chnl_ph": string(s.CommChnlPH), "new_chnl_ph": string(s.NewChnlPH), "cont_expr_ck": k}, e
	case ExprLabSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "lab", "comm_chnl_ph": string(s.CommChnlPH), "val_lab_qn": uniqsym.ConvertToString(s.ValLabQN), "cont_expr_ck": k}, e
	case ExprCaseSpec:
		branches := make(map[string]int64, len(s.ContExprs))
		for q, v := range s.ContExprs {
			k, e := child(v)
			if e != nil {
				return nil, e
			}
			name := uniqsym.ConvertToString(q)
			if _, ok := branches[name]; ok {
				return nil, fmt.Errorf("duplicate case label %q", name)
			}
			branches[name] = k
		}
		return map[string]any{"kind": "case", "comm_chnl_ph": string(s.CommChnlPH), "cont_expr_cks": branches}, nil
	case ExprLinkSpec:
		ys := make([]string, len(s.Ys))
		for i, v := range s.Ys {
			ys[i] = identity.ConvertToString(v)
		}
		return map[string]any{"kind": "link", "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "x": identity.ConvertToString(s.X), "ys": ys}, nil
	case ExprFwdSpec:
		return map[string]any{"kind": "fwd", "comm_chnl_ph": string(s.CommChnlPH), "cont_chnl_ph": string(s.ContChnlPH)}, nil
	case ExprCallSpec:
		k, e := child(s.ContExp)
		chs := make([]string, len(s.ValChnlPHs))
		for i, v := range s.ValChnlPHs {
			chs[i] = string(v)
		}
		return map[string]any{"kind": "call", "comm_chnl_ph": string(s.CommChnlPH), "new_chnl_ph": string(s.NewChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "val_chnl_phs": chs, "cont_expr_ck": k}, e
	case ExprSpawnSpec:
		k, e := child(s.ContExp)
		chs := make([]string, len(s.NewChnlPHs))
		for i, v := range s.NewChnlPHs {
			chs[i] = string(v)
		}
		return map[string]any{"kind": "spawn", "comm_chnl_ph": string(s.CommChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "new_chnl_phs": chs, "cont_expr_ck": k}, e
	case ExprAcquireSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "acquire", "comm_chnl_ph": string(s.CommChnlPH), "cont_expr_ck": k}, e
	case ExprAcceptSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "accept", "comm_chnl_ph": string(s.CommChnlPH), "cont_expr_ck": k}, e
	case ExprDetachSpec:
		return map[string]any{"kind": "detach", "comm_chnl_ph": string(s.CommChnlPH)}, nil
	case ExprReleaseSpec:
		return map[string]any{"kind": "release", "comm_chnl_ph": string(s.CommChnlPH)}, nil
	default:
		return nil, fmt.Errorf("unexpected process Definition spec %T", spec)
	}
}
