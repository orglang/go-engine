package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

// ContentKey returns the deterministic content key and canonical JSON
// representation of a process Term Definition expression.
func ContentKey(spec ExpSpec) (contkey.ADT, []byte, error) {
	value, err := canonicalExpSpec(spec)
	if err != nil {
		return contkey.Zero, nil, err
	}
	return contkey.HashJSON(value)
}

func canonicalExpSpec(spec ExpSpec) (any, error) {
	if spec == nil {
		return nil, fmt.Errorf("nil process term expression")
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
	case CloseSpec:
		return base("close", map[string]any{"cont_chnl_ph": string(s.ContChnlPH)}), nil
	case WaitSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("wait.cont: %w", err)
		}
		return base("wait", map[string]any{"cont_chnl_ph": string(s.ContChnlPH), "cont_exp": cont}), nil
	case SendSpec:
		return base("send", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "val_chnl_ph": string(s.ValChnlPH)}), nil
	case RecvSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("recv.cont: %w", err)
		}
		return base("recv", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "new_chnl_ph": string(s.NewChnlPH), "cont_exp": cont}), nil
	case LabSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("lab.cont: %w", err)
		}
		return base("lab", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "val_lab_qn": uniqsym.ConvertToString(s.ValLabQN), "cont_exp": cont}), nil
	case CaseSpec:
		branches := make(map[string]any, len(s.ContExps))
		for qn, child := range s.ContExps {
			name := uniqsym.ConvertToString(qn)
			if _, exists := branches[name]; exists {
				return nil, fmt.Errorf("duplicate case label %q", name)
			}
			value, err := nested(child)
			if err != nil {
				return nil, fmt.Errorf("case[%s]: %w", name, err)
			}
			branches[name] = value
		}
		return base("case", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exps": branches}), nil
	case LinkSpec:
		ys := make([]string, len(s.Ys))
		for i, id := range s.Ys {
			ys[i] = identity.ConvertToString(id)
		}
		return base("link", map[string]any{"proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "x": identity.ConvertToString(s.X), "ys": ys}), nil
	case FwdSpec:
		return base("fwd", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_chnl_ph": string(s.ContChnlPH)}), nil
	case CallSpec:
		chs := make([]string, len(s.ValChnlPHs))
		for i, ch := range s.ValChnlPHs {
			chs[i] = string(ch)
		}
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("call.cont: %w", err)
		}
		return base("call", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "new_chnl_ph": string(s.NewChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "val_chnl_phs": chs, "cont_exp": cont}), nil
	case SpawnSpec:
		chs := make([]string, len(s.NewChnlPHs))
		for i, ch := range s.NewChnlPHs {
			chs[i] = string(ch)
		}
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("spawn.cont: %w", err)
		}
		return base("spawn", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "new_chnl_phs": chs, "cont_exp": cont}), nil
	case AcqureSpec:
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
	case DetachSpec:
		return base("detach", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)}), nil
	case ReleaseSpec:
		return base("release", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)}), nil
	default:
		return nil, ErrExpTypeUnexpected(spec)
	}
}
