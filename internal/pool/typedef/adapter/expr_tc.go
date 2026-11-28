package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"

	typedef "orglang/go-engine/internal/pool/typedef/core"
)

func dataFromExprRef(ref typedef.ExpRef) ExpRef {
	if ref == nil {
		panic("can't be nil")
	}
	exprCK := contkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typedef.OneRef, typedef.OneRec:
		return ExpRef{K: OneKind, ExpCK: exprCK}
	case typedef.LinkRef, typedef.LinkRec:
		return ExpRef{K: LinkKind, ExpCK: exprCK}
	case typedef.PlusRef, typedef.PlusRec:
		return ExpRef{K: PlusKind, ExpCK: exprCK}
	case typedef.WithRef, typedef.WithRec:
		return ExpRef{K: WithKind, ExpCK: exprCK}
	default:
		panic(typedef.ErrRefTypeUnexpected(ref))
	}
}

func dataToExprRef(dto ExpRef) (typedef.ExpRef, error) {
	exprCK, err := contkey.ConvertFromInt(dto.ExpCK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case OneKind:
		return typedef.OneRef{ExpCK: exprCK}, nil
	case LinkKind:
		return typedef.LinkRef{ExpCK: exprCK}, nil
	case PlusKind:
		return typedef.PlusRef{ExpCK: exprCK}, nil
	case WithKind:
		return typedef.WithRef{ExpCK: exprCK}, nil
	default:
		panic(errExprKindUnexpected(dto.K))
	}
}

func dataToExprRec(dto ExpRec) (typedef.ExpRec, error) {
	states := make(map[int64]State, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ExpCK] = dto
	}
	return statesToExprRec(states, states[dto.ExpCK])
}

func dataFromExprRec(rec typedef.ExpRec) ExpRec {
	if rec == nil {
		panic("can't be nil")
	}
	dto := &ExpRec{
		ExpCK:  contkey.ConvertToInt(rec.Key()),
		States: nil,
	}
	statesFromExprRec(0, rec, dto)
	return *dto
}

func statesToExprRec(states map[int64]State, st State) (typedef.ExpRec, error) {
	exprCK, err := contkey.ConvertFromInt(st.ExpCK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case OneKind:
		return typedef.OneRec{ExpCK: exprCK}, nil
	case LinkKind:
		xactQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return typedef.LinkRec{ExpCK: exprCK, TypeQN: xactQN}, nil
	case PlusKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.Plus.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExpr, err := statesToExprRec(states, states[st.Spec.Plus.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.PlusRec{ExpCK: exprCK, ProcQNs: procQNs, ContExp: contExpr}, nil
	case WithKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.With.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExpr, err := statesToExprRec(states, states[st.Spec.With.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.WithRec{ExpCK: exprCK, ProcQNs: procQNs, ContExp: contExpr}, nil
	case UpKind:
		contExpr, err := statesToExprRec(states, states[st.Spec.Up.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.UpRec{ExpCK: exprCK, ContExp: contExpr}, nil
	case DownKind:
		contExpr, err := statesToExprRec(states, states[st.Spec.Down.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.DownRec{ExpCK: exprCK, ContExp: contExpr}, nil
	default:
		panic(errExprKindUnexpected(st.K))
	}
}

func statesFromExprRec(supExprCK int64, r typedef.ExpRec, dto *ExpRec) int64 {
	exprCK := contkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case typedef.OneRec:
		st := State{ExpCK: exprCK, K: OneKind, SupExprCK: supExprCK}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.LinkRec:
		st := State{
			ExpCK:    exprCK,
			K:        LinkKind,
			SupExprCK: supExprCK,
			Spec: ExpSpec{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.PlusRec:
		st := State{
			ExpCK:    exprCK,
			K:        PlusKind,
			SupExprCK: supExprCK,
			Spec: ExpSpec{Plus: &Sum{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpCK: statesFromExprRec(exprCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.WithRec:
		st := State{
			ExpCK:    exprCK,
			K:        WithKind,
			SupExprCK: supExprCK,
			Spec: ExpSpec{With: &Sum{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpCK: statesFromExprRec(exprCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.UpRec:
		st := State{
			ExpCK:    exprCK,
			K:        UpKind,
			SupExprCK: supExprCK,
			Spec: ExpSpec{Up: &Shift{
				ContExpCK: statesFromExprRec(exprCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.DownRec:
		st := State{
			ExpCK:    exprCK,
			K:        DownKind,
			SupExprCK: supExprCK,
			Spec: ExpSpec{Down: &Shift{
				ContExpCK: statesFromExprRec(exprCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return exprCK
	default:
		panic(typedef.ErrRecTypeUnexpected(r))
	}
}

func errExprKindUnexpected(k ExpKind) error {
	return fmt.Errorf("expr kind unexpected: %v", k)
}
