package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"

	typedef "orglang/go-engine/internal/proc/typedef/core"
)

func dataFromExprRef(ref typedef.ExpRef) ExpRef {
	if ref == nil {
		panic("can't be null")
	}
	expCK := contkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typedef.OneRef, typedef.OneRec:
		return ExpRef{K: OneKind, ExpCK: expCK}
	case typedef.LinkRef, typedef.LinkRec:
		return ExpRef{K: LinkKind, ExpCK: expCK}
	case typedef.TensorRef, typedef.TensorRec:
		return ExpRef{K: TensorKind, ExpCK: expCK}
	case typedef.LolliRef, typedef.LolliRec:
		return ExpRef{K: LolliKind, ExpCK: expCK}
	case typedef.PlusRef, typedef.PlusRec:
		return ExpRef{K: PlusKind, ExpCK: expCK}
	case typedef.WithRef, typedef.WithRec:
		return ExpRef{K: WithKind, ExpCK: expCK}
	default:
		panic(typedef.ErrRefTypeUnexpected(ref))
	}
}

func dataToExprRef(dto ExpRef) (typedef.ExpRef, error) {
	expCK, err := contkey.ConvertFromInt(dto.ExpCK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case OneKind:
		return typedef.OneRef{ExpCK: expCK}, nil
	case LinkKind:
		return typedef.LinkRef{ExpCK: expCK}, nil
	case TensorKind:
		return typedef.TensorRef{ExpCK: expCK}, nil
	case LolliKind:
		return typedef.LolliRef{ExpCK: expCK}, nil
	case PlusKind:
		return typedef.PlusRef{ExpCK: expCK}, nil
	case WithKind:
		return typedef.WithRef{ExpCK: expCK}, nil
	default:
		panic(errUnexpectedKind(dto.K))
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
	expCK, err := contkey.ConvertFromInt(st.ExpCK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case OneKind:
		return typedef.OneRec{ExpCK: expCK}, nil
	case LinkKind:
		typeQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return typedef.LinkRec{ExpCK: expCK, TypeQN: typeQN}, nil
	case TensorKind:
		b, err := statesToExprRec(states, states[st.Spec.Tensor.ValExpCK])
		if err != nil {
			return nil, err
		}
		c, err := statesToExprRec(states, states[st.Spec.Tensor.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.TensorRec{ExpCK: expCK, Val: b, Cont: c}, nil
	case LolliKind:
		y, err := statesToExprRec(states, states[st.Spec.Lolli.ValExpCK])
		if err != nil {
			return nil, err
		}
		z, err := statesToExprRec(states, states[st.Spec.Lolli.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.LolliRec{ExpCK: expCK, Val: y, Cont: z}, nil
	case PlusKind:
		choices := make(map[uniqsym.ADT]typedef.ExpRec, len(st.Spec.Plus))
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
		return typedef.PlusRec{ExpCK: expCK, Choices: choices}, nil
	case WithKind:
		choices := make(map[uniqsym.ADT]typedef.ExpRec, len(st.Spec.With))
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
		return typedef.WithRec{ExpCK: expCK, Choices: choices}, nil
	default:
		panic(errUnexpectedKind(st.K))
	}
}

func statesFromExprRec(fromID int64, r typedef.ExpRec, dto *ExpRec) int64 {
	expCK := contkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case typedef.OneRec:
		st := State{ExpCK: expCK, K: OneKind, SupExprCK: fromID}
		dto.States = append(dto.States, st)
		return expCK
	case typedef.LinkRec:
		st := State{
			ExpCK:    expCK,
			K:        LinkKind,
			SupExprCK: fromID,
			Spec: ExpSpec{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typedef.TensorRec:
		val := statesFromExprRec(expCK, rec.Val, dto)
		cont := statesFromExprRec(expCK, rec.Cont, dto)
		st := State{
			ExpCK:    expCK,
			K:        TensorKind,
			SupExprCK: fromID,
			Spec: ExpSpec{
				Tensor: &Prod{ValExpCK: val, ContExpCK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typedef.LolliRec:
		val := statesFromExprRec(expCK, rec.Val, dto)
		cont := statesFromExprRec(expCK, rec.Cont, dto)
		st := State{
			ExpCK:    expCK,
			K:        LolliKind,
			SupExprCK: fromID,
			Spec: ExpSpec{
				Lolli: &Prod{ValExpCK: val, ContExpCK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typedef.PlusRec:
		var choices []Sum
		for label, choice := range rec.Choices {
			cont := statesFromExprRec(expCK, choice, dto)
			choices = append(choices, Sum{LabQN: uniqsym.ConvertToString(label), ContExpCK: cont})
		}
		st := State{
			ExpCK:    expCK,
			K:        PlusKind,
			SupExprCK: fromID,
			Spec:     ExpSpec{Plus: choices},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typedef.WithRec:
		var choices []Sum
		for label, choice := range rec.Choices {
			cont := statesFromExprRec(expCK, choice, dto)
			choices = append(choices, Sum{LabQN: uniqsym.ConvertToString(label), ContExpCK: cont})
		}
		st := State{
			ExpCK:    expCK,
			K:        WithKind,
			SupExprCK: fromID,
			Spec:     ExpSpec{With: choices},
		}
		dto.States = append(dto.States, st)
		return expCK
	default:
		panic(typedef.ErrRecTypeUnexpected(r))
	}
}

func errUnexpectedKind(k ExpKind) error {
	return fmt.Errorf("unexpected kind %v", k)
}
