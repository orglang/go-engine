package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
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
	child := func(s ExprSpec) (int64, error) { k, _, e := ExprContentKey(s); return int64(k), e }
	switch s := spec.(type) {
	case ExprAcquireSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "acquire", "comm_chnl_ph": string(s.CommChnlPH), "cont_expr_ck": k}, e
	case ExprAcceptSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "accept", "comm_chnl_ph": string(s.CommChnlPH), "cont_expr_ck": k}, e
	case ExprHireSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "hire", "comm_chnl_ph": string(s.CommChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "cont_expr_ck": k}, e
	case ExprApplySpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "apply", "comm_chnl_ph": string(s.CommChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "cont_expr_ck": k}, e
	case ExprReleaseSpec:
		return map[string]any{"kind": "release", "comm_chnl_ph": string(s.CommChnlPH)}, nil
	case ExprDetachSpec:
		return map[string]any{"kind": "detach", "comm_chnl_ph": string(s.CommChnlPH)}, nil
	case ExprFireSpec:
		return map[string]any{"kind": "fire", "comm_chnl_ph": string(s.CommChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN)}, nil
	case ExprQuitSpec:
		return map[string]any{"kind": "quit", "comm_chnl_ph": string(s.CommChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN)}, nil
	case ExprSpawnSpec:
		qns := make([]string, len(s.ProcCompQNs))
		for i, q := range s.ProcCompQNs {
			qns[i] = uniqsym.ConvertToString(q)
		}
		return map[string]any{"kind": "spawn", "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "proc_comp_qns": qns}, nil
	default:
		return nil, fmt.Errorf("unexpected pool Definition spec %T", spec)
	}
}
