package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"
)

// ExpContentKey hashes only the declarative Definition model.
func ExpContentKey(spec ExpSpec) (contkey.ADT, []byte, error) {
	value, err := expContentValue(spec)
	if err != nil {
		return contkey.Zero, nil, err
	}
	return contkey.HashJSON(value)
}

func expContentValue(spec ExpSpec) (any, error) {
	child := func(s ExpSpec) (int64, error) { k, _, e := ExpContentKey(s); return int64(k), e }
	switch s := spec.(type) {
	case ExpAcquireSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "acquire", "comm_chnl_ph": string(s.CommChnlPH), "cont_exp_ck": k}, e
	case ExpAcceptSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "accept", "comm_chnl_ph": string(s.CommChnlPH), "cont_exp_ck": k}, e
	case ExpHireSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "hire", "comm_chnl_ph": string(s.CommChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "cont_exp_ck": k}, e
	case ExpApplySpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "apply", "comm_chnl_ph": string(s.CommChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "cont_exp_ck": k}, e
	case ExpReleaseSpec:
		return map[string]any{"kind": "release", "comm_chnl_ph": string(s.CommChnlPH)}, nil
	case ExpDetachSpec:
		return map[string]any{"kind": "detach", "comm_chnl_ph": string(s.CommChnlPH)}, nil
	case ExpFireSpec:
		return map[string]any{"kind": "fire", "comm_chnl_ph": string(s.CommChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN)}, nil
	case ExpQuitSpec:
		return map[string]any{"kind": "quit", "comm_chnl_ph": string(s.CommChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN)}, nil
	case ExpSpawnSpec:
		qns := make([]string, len(s.ProcCompQNs))
		for i, q := range s.ProcCompQNs {
			qns[i] = uniqsym.ConvertToString(q)
		}
		return map[string]any{"kind": "spawn", "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "proc_comp_qns": qns}, nil
	default:
		return nil, fmt.Errorf("unexpected pool Definition spec %T", spec)
	}
}
