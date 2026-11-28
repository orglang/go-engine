package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"

	typedef "orglang/go-engine/internal/proc/typedef/core"
)

func dataFromExprRef(ref typedef.ExprRef) ExprRef {
	if ref == nil {
		panic("can't be null")
	}
	exprCK := contkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typedef.OneRef, typedef.OneRec:
		return ExprRef{K: OneKind, ExprCK: exprCK}
	case typedef.LinkRef, typedef.LinkRec:
		return ExprRef{K: LinkKind, ExprCK: exprCK}
	case typedef.TensorRef, typedef.TensorRec:
		return ExprRef{K: TensorKind, ExprCK: exprCK}
	case typedef.LolliRef, typedef.LolliRec:
		return ExprRef{K: LolliKind, ExprCK: exprCK}
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
	case TensorKind:
		return typedef.TensorRef{ExprCK: exprCK}, nil
	case LolliKind:
		return typedef.LolliRef{ExprCK: exprCK}, nil
	case PlusKind:
		return typedef.PlusRef{ExprCK: exprCK}, nil
	case WithKind:
		return typedef.WithRef{ExprCK: exprCK}, nil
	default:
		panic(errUnexpectedKind(dto.K))
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
		typeQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return typedef.LinkRec{ExprCK: exprCK, TypeQN: typeQN}, nil
	case TensorKind:
		b, err := statesToExprRec(states, states[st.Spec.Tensor.ValExpCK])
		if err != nil {
			return nil, err
		}
		c, err := statesToExprRec(states, states[st.Spec.Tensor.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.TensorRec{ExprCK: exprCK, Val: b, Cont: c}, nil
	case LolliKind:
		y, err := statesToExprRec(states, states[st.Spec.Lolli.ValExpCK])
		if err != nil {
			return nil, err
		}
		z, err := statesToExprRec(states, states[st.Spec.Lolli.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.LolliRec{ExprCK: exprCK, Val: y, Cont: z}, nil
	case PlusKind:
		choices := make(map[uniqsym.ADT]typedef.ExprRec, len(st.Spec.Plus))
		for _, ch := range st.Spec.Plus {
			choice, err := statesToExprRec(states, states[ch.ContExpCK])
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(ch.LabQN)
			if err != nil {
				return nil, err
			}
			choices[label] = choice
		}
		return typedef.PlusRec{ExprCK: exprCK, Choices: choices}, nil
	case WithKind:
		choices := make(map[uniqsym.ADT]typedef.ExprRec, len(st.Spec.With))
		for _, ch := range st.Spec.With {
			choice, err := statesToExprRec(states, states[ch.ContExpCK])
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(ch.LabQN)
			if err != nil {
				return nil, err
			}
			choices[label] = choice
		}
		return typedef.WithRec{ExprCK: exprCK, Choices: choices}, nil
	default:
		panic(errUnexpectedKind(st.K))
	}
}

func statesFromExprRec(fromID int64, r typedef.ExprRec, dto *ExprRec) int64 {
	exprCK := contkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case typedef.OneRec:
		st := State{ExprCK: exprCK, K: OneKind, SupExprCK: fromID}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.LinkRec:
		st := State{
			ExprCK:     exprCK,
			K:         LinkKind,
			SupExprCK: fromID,
			Spec: ExprSpec{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.TensorRec:
		val := statesFromExprRec(exprCK, rec.Val, dto)
		cont := statesFromExprRec(exprCK, rec.Cont, dto)
		st := State{
			ExprCK:     exprCK,
			K:         TensorKind,
			SupExprCK: fromID,
			Spec: ExprSpec{
				Tensor: &Prod{ValExpCK: val, ContExpCK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.LolliRec:
		val := statesFromExprRec(exprCK, rec.Val, dto)
		cont := statesFromExprRec(exprCK, rec.Cont, dto)
		st := State{
			ExprCK:     exprCK,
			K:         LolliKind,
			SupExprCK: fromID,
			Spec: ExprSpec{
				Lolli: &Prod{ValExpCK: val, ContExpCK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.PlusRec:
		var choices []Sum
		for label, choice := range rec.Choices {
			cont := statesFromExprRec(exprCK, choice, dto)
			choices = append(choices, Sum{LabQN: uniqsym.ConvertToString(label), ContExpCK: cont})
		}
		st := State{
			ExprCK:     exprCK,
			K:         PlusKind,
			SupExprCK: fromID,
			Spec:      ExprSpec{Plus: choices},
		}
		dto.States = append(dto.States, st)
		return exprCK
	case typedef.WithRec:
		var choices []Sum
		for label, choice := range rec.Choices {
			cont := statesFromExprRec(exprCK, choice, dto)
			choices = append(choices, Sum{LabQN: uniqsym.ConvertToString(label), ContExpCK: cont})
		}
		st := State{
			ExprCK:     exprCK,
			K:         WithKind,
			SupExprCK: fromID,
			Spec:      ExprSpec{With: choices},
		}
		dto.States = append(dto.States, st)
		return exprCK
	default:
		panic(typedef.ErrRecTypeUnexpected(r))
	}
}

func errUnexpectedKind(k ExprKind) error {
	return fmt.Errorf("unexpected kind %v", k)
}
