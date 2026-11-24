package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
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
	child := func(s DescSpec) (int64, error) {
		k, _, err := DescContentKey(s)
		return int64(k), err
	}
	switch s := spec.(type) {
	case DescCloseSpec:
		return map[string]any{"kind": "close", "cont_chnl_ph": string(s.ContChnlPH)}, nil
	case DescWaitSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "wait", "cont_chnl_ph": string(s.ContChnlPH), "cont_exp_ck": k}, e
	case DescSendSpec:
		return map[string]any{"kind": "send", "conn_chnl_ph": string(s.ConnChnlPH), "val_chnl_ph": string(s.ValChnlPH)}, nil
	case DescRecvSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "recv", "conn_chnl_ph": string(s.ConnChnlPH), "new_chnl_ph": string(s.NewChnlPH), "cont_exp_ck": k}, e
	case DescLabSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "lab", "conn_chnl_ph": string(s.ConnChnlPH), "val_lab_qn": uniqsym.ConvertToString(s.ValLabQN), "cont_exp_ck": k}, e
	case DescCaseSpec:
		branches := make(map[string]int64, len(s.ContExps))
		for q, v := range s.ContExps {
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
		return map[string]any{"kind": "case", "conn_chnl_ph": string(s.ConnChnlPH), "cont_exp_cks": branches}, nil
	case DescLinkSpec:
		ys := make([]string, len(s.Ys))
		for i, v := range s.Ys {
			ys[i] = identity.ConvertToString(v)
		}
		return map[string]any{"kind": "link", "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "x": identity.ConvertToString(s.X), "ys": ys}, nil
	case DescFwdSpec:
		return map[string]any{"kind": "fwd", "conn_chnl_ph": string(s.ConnChnlPH), "cont_chnl_ph": string(s.ContChnlPH)}, nil
	case DescCallSpec:
		k, e := child(s.ContExp)
		chs := make([]string, len(s.ValChnlPHs))
		for i, v := range s.ValChnlPHs {
			chs[i] = string(v)
		}
		return map[string]any{"kind": "call", "conn_chnl_ph": string(s.ConnChnlPH), "new_chnl_ph": string(s.NewChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "val_chnl_phs": chs, "cont_exp_ck": k}, e
	case DescSpawnSpec:
		k, e := child(s.ContExp)
		chs := make([]string, len(s.NewChnlPHs))
		for i, v := range s.NewChnlPHs {
			chs[i] = string(v)
		}
		return map[string]any{"kind": "spawn", "conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "new_chnl_phs": chs, "cont_exp_ck": k}, e
	case DescAcquireSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "acquire", "conn_chnl_ph": string(s.ConnChnlPH), "cont_exp_ck": k}, e
	case DescAcceptSpec:
		k, e := child(s.ContExp)
		return map[string]any{"kind": "accept", "conn_chnl_ph": string(s.ConnChnlPH), "cont_exp_ck": k}, e
	case DescDetachSpec:
		return map[string]any{"kind": "detach", "conn_chnl_ph": string(s.ConnChnlPH)}, nil
	case DescReleaseSpec:
		return map[string]any{"kind": "release", "conn_chnl_ph": string(s.ConnChnlPH)}, nil
	default:
		return nil, fmt.Errorf("unexpected process Definition spec %T", spec)
	}
}
