package core

import (
	"orglang/go-engine/adt/valkey"
)

func ConvertSpecToRec(s ExpSpec) (ExpRec, error) {
	if s == nil {
		return nil, nil
	}
	switch spec := s.(type) {
	case OneSpec:
		return OneRec{ExpVK: valkey.One}, nil
	case LinkSpec:
		expVK, err := spec.TypeQN.Key()
		if err != nil {
			return nil, err
		}
		return LinkRec{ExpVK: expVK, TypeQN: spec.TypeQN}, nil
	case WithSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		vks := make([]valkey.ADT, 0, len(spec.ProcQNs)+1)
		for _, qn := range spec.ProcQNs {
			vk, err := qn.Key()
			if err != nil {
				return nil, err
			}
			vks = append(vks, vk)
		}
		expVK, err := valkey.Compose(append(vks, contExp.Key())...)
		if err != nil {
			return nil, err
		}
		return WithRec{ExpVK: expVK, ProcQNs: spec.ProcQNs, ContExp: contExp}, nil
	case PlusSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		vks := make([]valkey.ADT, 0, len(spec.ProcQNs)+1)
		for _, qn := range spec.ProcQNs {
			vk, err := qn.Key()
			if err != nil {
				return nil, err
			}
			vks = append(vks, vk)
		}
		expVK, err := valkey.Compose(append(vks, contExp.Key())...)
		return PlusRec{ExpVK: expVK, ProcQNs: spec.ProcQNs, ContExp: contExp}, nil
	case UpSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		expVK, err := valkey.Compose([]valkey.ADT{valkey.Two, contExp.Key()}...)
		return UpRec{ExpVK: expVK, ContExp: contExp}, nil
	case DownSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		expVK, err := valkey.Compose([]valkey.ADT{valkey.Three, contExp.Key()}...)
		return UpRec{ExpVK: expVK, ContExp: contExp}, nil
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
