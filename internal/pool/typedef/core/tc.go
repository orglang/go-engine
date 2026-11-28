package core

import (
	"orglang/go-engine/internal/adt/contkey"
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
		exprCK, err := contkey.Compose(append(cks, contExp.Key())...)
		if err != nil {
			return nil, err
		}
		return WithRec{ExprCK: exprCK, ProcQNs: spec.ProcQNs, ContExp: contExp}, nil
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
		exprCK, err := contkey.Compose(append(cks, contExp.Key())...)
		if err != nil {
			return nil, err
		}
		return PlusRec{ExprCK: exprCK, ProcQNs: spec.ProcQNs, ContExp: contExp}, nil
	case UpSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		exprCK, err := contkey.Compose([]contkey.ADT{contkey.Two, contExp.Key()}...)
		if err != nil {
			return nil, err
		}
		return UpRec{ExprCK: exprCK, ContExp: contExp}, nil
	case DownSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		exprCK, err := contkey.Compose([]contkey.ADT{contkey.Three, contExp.Key()}...)
		if err != nil {
			return nil, err
		}
		return DownRec{ExprCK: exprCK, ContExp: contExp}, nil
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
