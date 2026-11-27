package core

import (
	"orglang/go-engine/internal/adt/contkey"
)

func ConvertSpecToRec(s ExpSpec) (ExpRec, error) {
	if s == nil {
		return nil, nil
	}
	switch spec := s.(type) {
	case OneSpec:
		return OneRec{ExpCK: contkey.One}, nil
	case LinkSpec:
		expCK, err := spec.TypeQN.Key()
		if err != nil {
			return nil, err
		}
		return LinkRec{ExpCK: expCK, TypeQN: spec.TypeQN}, nil
	case WithSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		cks := make([]contkey.ADT, 0, len(spec.ProcQNs)+1)
		for _, qn := range spec.ProcQNs {
			ck, err := qn.Key()
			if err != nil {
				return nil, err
			}
			cks = append(cks, ck)
		}
		expCK, err := contkey.Compose(append(cks, contExp.Key())...)
		if err != nil {
			return nil, err
		}
		return WithRec{ExpCK: expCK, ProcQNs: spec.ProcQNs, ContExp: contExp}, nil
	case PlusSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		cks := make([]contkey.ADT, 0, len(spec.ProcQNs)+1)
		for _, qn := range spec.ProcQNs {
			ck, err := qn.Key()
			if err != nil {
				return nil, err
			}
			cks = append(cks, ck)
		}
		expCK, err := contkey.Compose(append(cks, contExp.Key())...)
		if err != nil {
			return nil, err
		}
		return PlusRec{ExpCK: expCK, ProcQNs: spec.ProcQNs, ContExp: contExp}, nil
	case UpSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		expCK, err := contkey.Compose([]contkey.ADT{contkey.Two, contExp.Key()}...)
		if err != nil {
			return nil, err
		}
		return UpRec{ExpCK: expCK, ContExp: contExp}, nil
	case DownSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		expCK, err := contkey.Compose([]contkey.ADT{contkey.Three, contExp.Key()}...)
		if err != nil {
			return nil, err
		}
		return DownRec{ExpCK: expCK, ContExp: contExp}, nil
	default:
		panic(ErrSpecTypeUnexpected(spec))
	}
}

func ConvertRecToSpec(r ExpRec) ExpSpec {
	if r == nil {
		return nil
	}
	switch rec := r.(type) {
	case OneRec:
		return OneSpec{}
	case LinkRec:
		return LinkSpec{TypeQN: rec.TypeQN}
	case WithRec:
		return WithSpec{ProcQNs: rec.ProcQNs, ContExp: ConvertRecToSpec(rec.ContExp)}
	case PlusRec:
		return PlusSpec{ProcQNs: rec.ProcQNs, ContExp: ConvertRecToSpec(rec.ContExp)}
	case UpRec:
		return UpSpec{ContExp: ConvertRecToSpec(rec.ContExp)}
	case DownRec:
		return DownSpec{ContExp: ConvertRecToSpec(rec.ContExp)}
	default:
		panic(ErrRecTypeUnexpected(rec))
	}
}
