package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"

	typedef "orglang/go-engine/internal/pool/typedef/core"
)

func dataFromExprRef(ref typedef.ExprRef) ExprRef {
	if ref == nil {
		panic("can't be nil")
	}
	exprCK := contkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typedef.OneRef, typedef.OneRec:
		return ExprRef{K: OneKind, ExprCK: exprCK}
	case typedef.LinkRef, typedef.LinkRec:
		return ExprRef{K: LinkKind, ExprCK: exprCK}
	case typedef.PlusRef, typedef.PlusRec:
		return ExprRef{K: PlusKind, ExprCK: exprCK}
	case typedef.WithRef, typedef.WithRec:
		return ExprRef{K: WithKind, ExprCK: exprCK}
	default:
		panic(typedef.ErrRefTypeUnexpected(ref))
	}
}

func dataToExprRef(dto ExprRef) (typedef.ExprRef, error) {
	exprCK, err := contkey.ConvertFromInt(dto.ExprCK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case OneKind:
		return typedef.OneRef{ExprCK: exprCK}, nil
	case LinkKind:
		return typedef.LinkRef{ExprCK: exprCK}, nil
	case PlusKind:
		return typedef.PlusRef{ExprCK: exprCK}, nil
	case WithKind:
		return typedef.WithRef{ExprCK: exprCK}, nil
	default:
		panic(errExprKindUnexpected(dto.K))
	}
}

func dataToExprRec(dto ExprRec) (typedef.ExprRec, error) {
	states := make(map[int64]State, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ExprCK] = dto
	}
	return statesToExprRec(states, states[dto.ExprCK])
}

func dataFromExprRec(rec typedef.ExprRec) ExprRec {
	if rec == nil {
		panic("can't be nil")
	}
	dto := &ExprRec{
		ExprCK:  contkey.ConvertToInt(rec.Key()),
		States: nil,
	}
	statesFromExprRec(0, rec, dto)
	return *dto
}

func statesToExprRec(states map[int64]State, st State) (typedef.ExprRec, error) {
	exprCK, err := contkey.ConvertFromInt(st.ExprCK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case OneKind:
		return typedef.OneRec{ExprCK: exprCK}, nil
	case LinkKind:
		xactQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return typedef.LinkRec{ExprCK: exprCK, TypeQN: xactQN}, nil
	case PlusKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.Plus.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExpr, err := statesToExprRec(states, states[st.Spec.Plus.ContExprCK])
		if err != nil {
			return nil, err
		}
		return typedef.PlusRec{ExprCK: exprCK, ProcQNs: procQNs, ContExpr: contExpr}, nil
	case WithKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.With.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExpr, err := statesToExprRec(states, states[st.Spec.With.ContExprCK])
		if err != nil {
			return nil, err
		}
		return typedef.WithRec{ExprCK: exprCK, ProcQNs: procQNs, ContExpr: contExpr}, nil
	case UpKind:
		contExpr, err := statesToExprRec(states, states[st.Spec.Up.ContExprCK])
		if err != nil {
			return nil, err
		}
		return typedef.UpRec{ExprCK: exprCK, ContExpr: contExpr}, nil
	case DownKind:
		contExpr, err := statesToExprRec(states, states[st.Spec.Down.ContExprCK])
		if err != nil {
			return nil, err
		}
		return typedef.DownRec{ExprCK: exprCK, ContExpr: contExpr}, nil
	default:
		panic(errExprKindUnexpected(st.K))
	}
}

func statesFromExprRec(supExprCK int64, r typedef.ExprRec, dto *ExprRec) int64 {
	exprCK := contkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case typedef.OneRec:
		st := State{ExprCK: exprCK, K: OneKind, SupExprCK: supExprCK}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.LinkRec:
		st := State{
			ExprCK:     exprCK,
			K:         LinkKind,
			SupExprCK: supExprCK,
			Spec: ExprSpec{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.PlusRec:
		st := State{
			ExprCK:     exprCK,
			K:         PlusKind,
			SupExprCK: supExprCK,
			Spec: ExprSpec{Plus: &Sum{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExprCK: statesFromExprRec(exprCK, rec.ContExpr, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.WithRec:
		st := State{
			ExprCK:     exprCK,
			K:         WithKind,
			SupExprCK: supExprCK,
			Spec: ExprSpec{With: &Sum{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExprCK: statesFromExprRec(exprCK, rec.ContExpr, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.UpRec:
		st := State{
			ExprCK:     exprCK,
			K:         UpKind,
			SupExprCK: supExprCK,
			Spec: ExprSpec{Up: &Shift{
				ContExprCK: statesFromExprRec(exprCK, rec.ContExpr, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.DownRec:
		st := State{
			ExprCK:     exprCK,
			K:         DownKind,
			SupExprCK: supExprCK,
			Spec: ExprSpec{Down: &Shift{
				ContExprCK: statesFromExprRec(exprCK, rec.ContExpr, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return exprCK
	default:
		panic(typedef.ErrRecTypeUnexpected(r))
	}
}

func errExprKindUnexpected(k ExprKind) error {
	return fmt.Errorf("expr kind unexpected: %v", k)
}
