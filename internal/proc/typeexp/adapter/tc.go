package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"

	typeexp "orglang/go-engine/internal/proc/typeexp/core"
)

func dataFromExpRef(ref typeexp.ExpRef) ExpRef {
	if ref == nil {
		panic("can't be null")
	}
	expCK := contkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typeexp.OneRef, typeexp.OneRec:
		return ExpRef{K: OneKind, ExpCK: expCK}
	case typeexp.LinkRef, typeexp.LinkRec:
		return ExpRef{K: LinkKind, ExpCK: expCK}
	case typeexp.TensorRef, typeexp.TensorRec:
		return ExpRef{K: TensorKind, ExpCK: expCK}
	case typeexp.LolliRef, typeexp.LolliRec:
		return ExpRef{K: LolliKind, ExpCK: expCK}
	case typeexp.PlusRef, typeexp.PlusRec:
		return ExpRef{K: PlusKind, ExpCK: expCK}
	case typeexp.WithRef, typeexp.WithRec:
		return ExpRef{K: WithKind, ExpCK: expCK}
	default:
		panic(typeexp.ErrRefTypeUnexpected(ref))
	}
}

func dataToExpRef(dto ExpRef) (typeexp.ExpRef, error) {
	expCK, err := contkey.ConvertFromInt(dto.ExpCK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case OneKind:
		return typeexp.OneRef{ExpCK: expCK}, nil
	case LinkKind:
		return typeexp.LinkRef{ExpCK: expCK}, nil
	case TensorKind:
		return typeexp.TensorRef{ExpCK: expCK}, nil
	case LolliKind:
		return typeexp.LolliRef{ExpCK: expCK}, nil
	case PlusKind:
		return typeexp.PlusRef{ExpCK: expCK}, nil
	case WithKind:
		return typeexp.WithRef{ExpCK: expCK}, nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func dataToExpRec(dto ExpRec) (typeexp.ExpRec, error) {
	states := make(map[int64]State, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ExpCK] = dto
	}
	return statesToExpRec(states, states[dto.ExpCK])
}

func dataFromExpRec(rec typeexp.ExpRec) ExpRec {
	if rec == nil {
		panic("can't be nil")
	}
	dto := &ExpRec{
		ExpCK:  contkey.ConvertToInt(rec.Key()),
		States: nil,
	}
	statesFromExpRec(0, rec, dto)
	return *dto
}

func statesToExpRec(states map[int64]State, st State) (typeexp.ExpRec, error) {
	expCK, err := contkey.ConvertFromInt(st.ExpCK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case OneKind:
		return typeexp.OneRec{ExpCK: expCK}, nil
	case LinkKind:
		typeQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return typeexp.LinkRec{ExpCK: expCK, TypeQN: typeQN}, nil
	case TensorKind:
		b, err := statesToExpRec(states, states[st.Spec.Tensor.ValExpCK])
		if err != nil {
			return nil, err
		}
		c, err := statesToExpRec(states, states[st.Spec.Tensor.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typeexp.TensorRec{ExpCK: expCK, Val: b, Cont: c}, nil
	case LolliKind:
		y, err := statesToExpRec(states, states[st.Spec.Lolli.ValExpCK])
		if err != nil {
			return nil, err
		}
		z, err := statesToExpRec(states, states[st.Spec.Lolli.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typeexp.LolliRec{ExpCK: expCK, Val: y, Cont: z}, nil
	case PlusKind:
		choices := make(map[uniqsym.ADT]typeexp.ExpRec, len(st.Spec.Plus))
		for _, ch := range st.Spec.Plus {
			choice, err := statesToExpRec(states, states[ch.ContExpCK])
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(ch.LabQN)
			if err != nil {
				return nil, err
			}
			choices[label] = choice
		}
		return typeexp.PlusRec{ExpCK: expCK, Choices: choices}, nil
	case WithKind:
		choices := make(map[uniqsym.ADT]typeexp.ExpRec, len(st.Spec.With))
		for _, ch := range st.Spec.With {
			choice, err := statesToExpRec(states, states[ch.ContExpCK])
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(ch.LabQN)
			if err != nil {
				return nil, err
			}
			choices[label] = choice
		}
		return typeexp.WithRec{ExpCK: expCK, Choices: choices}, nil
	default:
		panic(errUnexpectedKind(st.K))
	}
}

func statesFromExpRec(fromID int64, r typeexp.ExpRec, dto *ExpRec) int64 {
	expCK := contkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case typeexp.OneRec:
		st := State{ExpCK: expCK, K: OneKind, SupExpCK: fromID}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.LinkRec:
		st := State{
			ExpCK:    expCK,
			K:        LinkKind,
			SupExpCK: fromID,
			Spec: ExpSpec{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.TensorRec:
		val := statesFromExpRec(expCK, rec.Val, dto)
		cont := statesFromExpRec(expCK, rec.Cont, dto)
		st := State{
			ExpCK:    expCK,
			K:        TensorKind,
			SupExpCK: fromID,
			Spec: ExpSpec{
				Tensor: &Prod{ValExpCK: val, ContExpCK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.LolliRec:
		val := statesFromExpRec(expCK, rec.Val, dto)
		cont := statesFromExpRec(expCK, rec.Cont, dto)
		st := State{
			ExpCK:    expCK,
			K:        LolliKind,
			SupExpCK: fromID,
			Spec: ExpSpec{
				Lolli: &Prod{ValExpCK: val, ContExpCK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.PlusRec:
		var choices []Sum
		for label, choice := range rec.Choices {
			cont := statesFromExpRec(expCK, choice, dto)
			choices = append(choices, Sum{LabQN: uniqsym.ConvertToString(label), ContExpCK: cont})
		}
		st := State{
			ExpCK:    expCK,
			K:        PlusKind,
			SupExpCK: fromID,
			Spec:     ExpSpec{Plus: choices},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.WithRec:
		var choices []Sum
		for label, choice := range rec.Choices {
			cont := statesFromExpRec(expCK, choice, dto)
			choices = append(choices, Sum{LabQN: uniqsym.ConvertToString(label), ContExpCK: cont})
		}
		st := State{
			ExpCK:    expCK,
			K:        WithKind,
			SupExpCK: fromID,
			Spec:     ExpSpec{With: choices},
		}
		dto.States = append(dto.States, st)
		return expCK
	default:
		panic(typeexp.ErrRecTypeUnexpected(r))
	}
}

func errUnexpectedKind(k ExpKind) error {
	return fmt.Errorf("unexpected kind %v", k)
}
