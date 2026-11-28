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
		contExpr, err := ConvertSpecToRec(spec.ContExpr)
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
		exprCK, err := contkey.Compose(append(cks, contExpr.Key())...)
		if err != nil {
			return nil, err
		}
		return WithRec{ExprCK: exprCK, ProcQNs: spec.ProcQNs, ContExpr: contExpr}, nil
	case PlusSpec:
		contExpr, err := ConvertSpecToRec(spec.ContExpr)
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
		exprCK, err := contkey.Compose(append(cks, contExpr.Key())...)
		if err != nil {
			return nil, err
		}
		return PlusRec{ExprCK: exprCK, ProcQNs: spec.ProcQNs, ContExpr: contExpr}, nil
	case UpSpec:
		contExpr, err := ConvertSpecToRec(spec.ContExpr)
		if err != nil {
			return nil, err
		}
		exprCK, err := contkey.Compose([]contkey.ADT{contkey.Two, contExpr.Key()}...)
		if err != nil {
			return nil, err
		}
		return UpRec{ExprCK: exprCK, ContExpr: contExpr}, nil
	case DownSpec:
		contExpr, err := ConvertSpecToRec(spec.ContExpr)
		if err != nil {
			return nil, err
		}
		exprCK, err := contkey.Compose([]contkey.ADT{contkey.Three, contExpr.Key()}...)
		if err != nil {
			return nil, err
		}
		return DownRec{ExprCK: exprCK, ContExpr: contExpr}, nil
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
		return WithSpec{ProcQNs: rec.ProcQNs, ContExpr: ConvertRecToSpec(rec.ContExpr)}
	case PlusRec:
		return PlusSpec{ProcQNs: rec.ProcQNs, ContExpr: ConvertRecToSpec(rec.ContExpr)}
	case UpRec:
		return UpSpec{ContExpr: ConvertRecToSpec(rec.ContExpr)}
	case DownRec:
		return DownSpec{ContExpr: ConvertRecToSpec(rec.ContExpr)}
	default:
		panic(ErrRecTypeUnexpected(rec))
	}
}
