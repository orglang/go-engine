package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"

	termexp "orglang/go-engine/internal/pool/termexp/core"
)

func DataFromExpSpec(s termexp.ExpSpec) ExpSpec {
	switch spec := s.(type) {
	case termexp.AcquireSpec:
		return ExpSpec{K: AcquireKind, Acquire: DataFromAcquireSpec(spec)}
	case termexp.AcceptSpec:
		return ExpSpec{K: AcceptKind, Accept: DataFromAcceptSpec(spec)}
	case termexp.HireSpec:
		return ExpSpec{K: HireKind, Hire: DataFromHireSpec(spec)}
	case termexp.ApplySpec:
		return ExpSpec{K: ApplyKind, Apply: DataFromApplySpec(spec)}
	case termexp.ReleaseSpec:
		return ExpSpec{K: ReleaseKind, Release: DataFromReleaseSpec(spec)}
	case termexp.DetachSpec:
		return ExpSpec{K: DetachKind, Detach: DataFromDetachSpec(spec)}
	default:
		panic(termexp.ErrSpecTypeUnexpected(s))
	}
}

func DataToExpSpec(dto ExpSpec) (termexp.ExpSpec, error) {
	switch dto.K {
	case AcquireKind:
		return DataToAcquireSpec(dto.Acquire)
	case AcceptKind:
		return DataToAcceptSpec(dto.Accept)
	case HireKind:
		return DataToHireSpec(dto.Hire)
	case ApplyKind:
		return DataToApplySpec(dto.Apply)
	case ReleaseKind:
		return DataToReleaseSpec(dto.Release)
	case DetachKind:
		return DataToDetachSpec(dto.Detach)
	default:
		panic(fmt.Errorf("exp kind unexpected: %v", dto.K))
	}
}

func DataFromExpRec(r termexp.ExpRec) ExpRec {
	switch rec := r.(type) {
	case termexp.AcquireRec:
		return ExpRec{K: AcquireKind, Acquire: DataFromAcquireRec(rec)}
	case termexp.AcceptRec:
		return ExpRec{K: AcceptKind, Accept: DataFromAcceptRec(rec)}
	case termexp.HireRec:
		return ExpRec{K: HireKind, Hire: DataFromHireRec(rec)}
	case termexp.ApplyRec:
		return ExpRec{K: ApplyKind, Apply: DataFromApplyRec(rec)}
	case termexp.ReleaseRec:
		return ExpRec{K: ReleaseKind, Release: DataFromReleaseRec(rec)}
	case termexp.DetachRec:
		return ExpRec{K: DetachKind, Detach: DataFromDetachRec(rec)}
	default:
		panic(termexp.ErrRecTypeUnexpected(r))
	}
}

func DataToExpRec(dto ExpRec) (termexp.ExpRec, error) {
	switch dto.K {
	case AcquireKind:
		return DataToAcquireRec(dto.Acquire)
	case AcceptKind:
		return DataToAcceptRec(dto.Accept)
	case HireKind:
		return DataToHireRec(dto.Hire)
	case ApplyKind:
		return DataToApplyRec(dto.Apply)
	case ReleaseKind:
		return DataToReleaseRec(dto.Release)
	case DetachKind:
		return DataToDetachRec(dto.Detach)
	default:
		panic(fmt.Errorf("exp kind unexpected: %v", dto.K))
	}
}

// canonicalExpSpec converts a domain expression into the stable JSONB storage shape.
func canonicalExpSpec(spec ExpSpec) (any, error) {
	if spec == nil {
		return nil, fmt.Errorf("nil pool term expression")
	}

	base := func(kind string, fields map[string]any) map[string]any {
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
