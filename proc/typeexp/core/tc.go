package core

import (
	"orglang/go-engine/adt/uniqsym"
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
	case TensorSpec:
		val, err := ConvertSpecToRec(spec.Val)
		if err != nil {
			return nil, err
		}
		cont, err := ConvertSpecToRec(spec.Cont)
		if err != nil {
			return nil, err
		}
		key := val.Key() + cont.Key()
		return TensorRec{ExpVK: key, Val: val, Cont: cont}, nil
	case LolliSpec:
		val, err := ConvertSpecToRec(spec.Val)
		if err != nil {
			return nil, err
		}
		cont, err := ConvertSpecToRec(spec.Cont)
		if err != nil {
			return nil, err
		}
		key := val.Key() + cont.Key()
		return LolliRec{ExpVK: key, Val: val, Cont: cont}, nil
	case WithSpec:
		conts := make(map[uniqsym.ADT]ExpRec, len(spec.Choices))
		keys := make([]valkey.ADT, len(spec.Choices)*2)
		for lab, spec := range spec.Choices {
			cont, err := ConvertSpecToRec(spec)
			if err != nil {
				return nil, err
			}
			conts[lab] = cont
			labVK, err := lab.Key()
			if err != nil {
				return nil, err
			}
			keys = append(keys, labVK, cont.Key())
		}
		expVK, err := valkey.Compose(keys...)
		if err != nil {
			return nil, err
		}
		return WithRec{ExpVK: expVK, Choices: conts}, nil
	case PlusSpec:
		conts := make(map[uniqsym.ADT]ExpRec, len(spec.Choices))
		keys := make([]valkey.ADT, len(spec.Choices)*2)
		for lab, spec := range spec.Choices {
			cont, err := ConvertSpecToRec(spec)
			if err != nil {
				return nil, err
			}
			conts[lab] = cont
			labVK, err := lab.Key()
			if err != nil {
				return nil, err
			}
			keys = append(keys, labVK, cont.Key())
		}
		expVK, err := valkey.Compose(keys...)
		if err != nil {
			return nil, err
		}
		return PlusRec{ExpVK: expVK, Choices: conts}, nil
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
	case TensorRec:
		return TensorSpec{
			Val:  ConvertRecToSpec(rec.Val),
			Cont: ConvertRecToSpec(rec.Cont),
		}
	case LolliRec:
		return LolliSpec{
			Val:  ConvertRecToSpec(rec.Val),
			Cont: ConvertRecToSpec(rec.Cont),
		}
	case WithRec:
		choices := make(map[uniqsym.ADT]ExpSpec, len(rec.Choices))
		for lab, cont := range rec.Choices {
			choices[lab] = ConvertRecToSpec(cont)
		}
		return WithSpec{Choices: choices}
	case PlusRec:
		choices := make(map[uniqsym.ADT]ExpSpec, len(rec.Choices))
		for lab, cont := range rec.Choices {
			choices[lab] = ConvertRecToSpec(cont)
		}
		return PlusSpec{Choices: choices}
	default:
		panic(ErrRecTypeUnexpected(rec))
	}
}
