package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"
)

// DescContentKey hashes only the declarative Definition model.
func DescContentKey(spec DescSpec) (contkey.ADT, []byte, error) {
	value, err := descContentValue(spec)
	if err != nil {
		return contkey.Zero, nil, err
	}
	return contkey.HashJSON(value)
}

func descContentValue(spec DescSpec) (any, error) {
	child := func(s DescSpec) (int64, error) { k, _, e := DescContentKey(s); return int64(k), e }
	switch s := spec.(type) {
	case DescAcquireSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "acquire", "conn_chnl_ph": string(s.ConnChnlPH), "cont_exp_ck": k}, e
	case DescAcceptSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "accept", "conn_chnl_ph": string(s.ConnChnlPH), "cont_exp_ck": k}, e
	case DescHireSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "hire", "conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "cont_exp_ck": k}, e
	case DescApplySpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "apply", "conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "cont_exp_ck": k}, e
	case DescReleaseSpec:
		return map[string]any{"kind": "release", "conn_chnl_ph": string(s.ConnChnlPH)}, nil
	case DescDetachSpec:
		return map[string]any{"kind": "detach", "conn_chnl_ph": string(s.ConnChnlPH)}, nil
	case DescFireSpec:
		return map[string]any{"kind": "fire", "conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN)}, nil
	case DescQuitSpec:
		return map[string]any{"kind": "quit", "conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN)}, nil
	case DescSpawnSpec:
		qns := make([]string, len(s.ProcCompQNs))
		for i, q := range s.ProcCompQNs {
			qns[i] = uniqsym.ConvertToString(q)
		}
		return map[string]any{"kind": "spawn", "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "proc_comp_qns": qns}, nil
	default:
		return nil, fmt.Errorf("unexpected pool Definition spec %T", spec)
	}
}
