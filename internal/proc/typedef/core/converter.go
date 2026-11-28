package core

import (
	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"
)

func ConvertSpecToRec(s ExprSpec) (ExprRec, error) {
	if s == nil {
		return nil, nil
	}
	switch spec := s.(type) {
	case OneSpec:
		return OneRec{ExprCK: contkey.One}, nil
	case LinkSpec:
		exprCK, err := spec.TypeQN.Key()
		if err != nil {
			return nil, err
		}
		return LinkRec{ExprCK: exprCK, TypeQN: spec.TypeQN}, nil
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
		return TensorRec{ExprCK: key, Val: val, Cont: cont}, nil
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
		return LolliRec{ExprCK: key, Val: val, Cont: cont}, nil
	case WithSpec:
		conts := make(map[uniqsym.ADT]ExprRec, len(spec.Choices))
		keys := make([]contkey.ADT, len(spec.Choices)*2)
		for lab, spec := range spec.Choices {
			cont, err := ConvertSpecToRec(spec)
			if err != nil {
				return nil, err
			}
			conts[lab] = cont
			labCK, err := lab.Key()
			if err != nil {
				return nil, err
			}
			keys = append(keys, labCK, cont.Key())
		}
		exprCK, err := contkey.Compose(keys...)
		if err != nil {
			return nil, err
		}
		return WithRec{ExprCK: exprCK, Choices: conts}, nil
	case PlusSpec:
		conts := make(map[uniqsym.ADT]ExprRec, len(spec.Choices))
		keys := make([]contkey.ADT, len(spec.Choices)*2)
		for lab, spec := range spec.Choices {
			cont, err := ConvertSpecToRec(spec)
			if err != nil {
				return nil, err
			}
			conts[lab] = cont
			labCK, err := lab.Key()
			if err != nil {
				return nil, err
			}
			keys = append(keys, labCK, cont.Key())
		}
		exprCK, err := contkey.Compose(keys...)
		if err != nil {
			return nil, err
		}
		return PlusRec{ExprCK: exprCK, Choices: conts}, nil
	default:
		panic(ErrSpecTypeUnexpected(spec))
	}
}

func ConvertRecToSpec(r ExprRec) ExprSpec {
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
		choices := make(map[uniqsym.ADT]ExprSpec, len(rec.Choices))
		for lab, cont := range rec.Choices {
			choices[lab] = ConvertRecToSpec(cont)
		}
		return WithSpec{Choices: choices}
	case PlusRec:
		choices := make(map[uniqsym.ADT]ExprSpec, len(rec.Choices))
		for lab, cont := range rec.Choices {
			choices[lab] = ConvertRecToSpec(cont)
		}
		return PlusSpec{Choices: choices}
	default:
		panic(ErrRecTypeUnexpected(rec))
	}
}
