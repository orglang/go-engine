package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"
)

// ContentKey returns the deterministic content key and canonical JSON
// representation of a pool Term Definition expression.
func ContentKey(spec ExpSpec) (contkey.ADT, []byte, error) {
	value, err := canonicalExpSpec(spec)
	if err != nil {
		return contkey.Zero, nil, err
	}
	return contkey.HashJSON(value)
}

func canonicalExpSpec(spec ExpSpec) (any, error) {
	if spec == nil {
		return nil, fmt.Errorf("nil pool term expression")
	}

	base := func(kind string, fields map[string]any) map[string]any {
		fields["format"] = "expck-json-v1"
		fields["kind"] = kind
		return fields
	}
	nested := func(child ExpSpec) (any, error) {
		return canonicalExpSpec(child)
	}

	switch s := spec.(type) {
	case AcquireSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("acquire.cont: %w", err)
		}
		return base("acquire", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp": cont}), nil
	case AcceptSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("accept.cont: %w", err)
		}
		return base("accept", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp": cont}), nil
	case HireSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("hire.cont: %w", err)
		}
		return base("hire", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "cont_exp": cont}), nil
	case ApplySpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("apply.cont: %w", err)
		}
		return base("apply", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "cont_exp": cont}), nil
	case ReleaseSpec:
		return base("release", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)}), nil
	case DetachSpec:
		return base("detach", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)}), nil
	case FireSpec:
		return base("fire", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN)}), nil
	case QuitSpec:
		return base("quit", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN)}), nil
	case SpawnSpec:
		qns := make([]string, len(s.ProcCompQNs))
		for i, qn := range s.ProcCompQNs {
			qns[i] = uniqsym.ConvertToString(qn)
		}
		return base("spawn", map[string]any{"proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "proc_comp_qns": qns}), nil
	case SpawnSpec2:
		return nil, fmt.Errorf("pool SpawnSpec2 is an execution reference form and has no Definition content key")
	default:
		return nil, ErrSpecTypeUnexpected(spec)
	}
}
