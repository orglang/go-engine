package adapter

import (
	"fmt"

	"golang.org/x/exp/maps"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"

	sdktypeexp "github.com/orglang/go-sdk/proc/typeexp"

	typeexp "orglang/go-engine/internal/proc/typeexp/core"
)

func MsgFromExpSpec(s typeexp.ExpSpec) sdktypeexp.ExpSpec {
	switch spec := s.(type) {
	case typeexp.OneSpec:
		return sdktypeexp.ExpSpec{K: sdktypeexp.One}
	case typeexp.LinkSpec:
		return sdktypeexp.ExpSpec{
			K:    sdktypeexp.Link,
			Link: &sdktypeexp.LinkSpec{TypeQN: uniqsym.ConvertToString(spec.TypeQN)}}
	case typeexp.TensorSpec:
		return sdktypeexp.ExpSpec{
			K: sdktypeexp.Tensor,
			Tensor: &sdktypeexp.ProdSpec{
				Val:  MsgFromExpSpec(spec.Val),
				Cont: MsgFromExpSpec(spec.Cont),
			},
		}
	case typeexp.LolliSpec:
		return sdktypeexp.ExpSpec{
			K: sdktypeexp.Lolli,
			Lolli: &sdktypeexp.ProdSpec{
				Val:  MsgFromExpSpec(spec.Val),
				Cont: MsgFromExpSpec(spec.Cont),
			},
		}
	case typeexp.WithSpec:
		choices := make([]sdktypeexp.ChoiceSpec, len(spec.Choices))
		for i, l := range maps.Keys(spec.Choices) {
			choices[i] = sdktypeexp.ChoiceSpec{
				LabQN: uniqsym.ConvertToString(l),
				Cont:  MsgFromExpSpec(spec.Choices[l]),
			}
		}
		return sdktypeexp.ExpSpec{K: sdktypeexp.With, With: &sdktypeexp.SumSpec{Choices: choices}}
	case typeexp.PlusSpec:
		choices := make([]sdktypeexp.ChoiceSpec, len(spec.Choices))
		for i, l := range maps.Keys(spec.Choices) {
			choices[i] = sdktypeexp.ChoiceSpec{
				LabQN: uniqsym.ConvertToString(l),
				Cont:  MsgFromExpSpec(spec.Choices[l]),
			}
		}
		return sdktypeexp.ExpSpec{K: sdktypeexp.Plus, Plus: &sdktypeexp.SumSpec{Choices: choices}}
	default:
		panic(typeexp.ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpec(dto sdktypeexp.ExpSpec) (typeexp.ExpSpec, error) {
	switch dto.K {
	case sdktypeexp.One:
		return typeexp.OneSpec{}, nil
	case sdktypeexp.Link:
		typeQN, err := uniqsym.ConvertFromString(dto.Link.TypeQN)
		if err != nil {
			return nil, err
		}
		return typeexp.LinkSpec{TypeQN: typeQN}, nil
	case sdktypeexp.Tensor:
		valES, err := MsgToExpSpec(dto.Tensor.Val)
		if err != nil {
			return nil, err
		}
		contES, err := MsgToExpSpec(dto.Tensor.Cont)
		if err != nil {
			return nil, err
		}
		return typeexp.TensorSpec{Val: valES, Cont: contES}, nil
	case sdktypeexp.Lolli:
		valES, err := MsgToExpSpec(dto.Lolli.Val)
		if err != nil {
			return nil, err
		}
		contES, err := MsgToExpSpec(dto.Lolli.Cont)
		if err != nil {
			return nil, err
		}
		return typeexp.LolliSpec{Val: valES, Cont: contES}, nil
	case sdktypeexp.Plus:
		choices := make(map[uniqsym.ADT]typeexp.ExpSpec, len(dto.Plus.Choices))
		for _, ch := range dto.Plus.Choices {
			choice, err := MsgToExpSpec(ch.Cont)
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(ch.LabQN)
			if err != nil {
				return nil, err
			}
			choices[label] = choice
		}
		return typeexp.PlusSpec{Choices: choices}, nil
	case sdktypeexp.With:
		choices := make(map[uniqsym.ADT]typeexp.ExpSpec, len(dto.With.Choices))
		for _, ch := range dto.With.Choices {
			choice, err := MsgToExpSpec(ch.Cont)
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(ch.LabQN)
			if err != nil {
				return nil, err
			}
			choices[label] = choice
		}
		return typeexp.WithSpec{Choices: choices}, nil
	default:
		panic(sdktypeexp.ErrKindUnexpected(dto.K))
	}
}

func MsgFromExpRef(ref typeexp.ExpRef) sdktypeexp.ExpRef {
	expCK := contkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typeexp.OneRef, typeexp.OneRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.One, ExpVK: expCK}
	case typeexp.LinkRef, typeexp.LinkRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.Link, ExpVK: expCK}
	case typeexp.TensorRef, typeexp.TensorRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.Tensor, ExpVK: expCK}
	case typeexp.LolliRef, typeexp.LolliRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.Lolli, ExpVK: expCK}
	case typeexp.PlusRef, typeexp.PlusRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.Plus, ExpVK: expCK}
	case typeexp.WithRef, typeexp.WithRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.With, ExpVK: expCK}
	default:
		panic(typeexp.ErrRefTypeUnexpected(ref))
	}
}

func MsgToExpRef(dto sdktypeexp.ExpRef) (typeexp.ExpRef, error) {
	expCK, err := contkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case sdktypeexp.One:
		return typeexp.OneRef{ExpCK: expCK}, nil
	case sdktypeexp.Link:
		return typeexp.LinkRef{ExpCK: expCK}, nil
	case sdktypeexp.Tensor:
		return typeexp.TensorRef{ExpCK: expCK}, nil
	case sdktypeexp.Lolli:
		return typeexp.LolliRef{ExpCK: expCK}, nil
	case sdktypeexp.Plus:
		return typeexp.PlusRef{ExpCK: expCK}, nil
	case sdktypeexp.With:
		return typeexp.WithRef{ExpCK: expCK}, nil
	default:
		panic(sdktypeexp.ErrKindUnexpected(dto.K))
	}
}

func dataFromExpRef(ref typeexp.ExpRef) typeexp.ExpRefDS {
	if ref == nil {
		panic("can't be null")
	}
	expCK := contkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typeexp.OneRef, typeexp.OneRec:
		return typeexp.ExpRefDS{K: typeexp.OneKind, ExpCK: expCK}
	case typeexp.LinkRef, typeexp.LinkRec:
		return typeexp.ExpRefDS{K: typeexp.LinkKind, ExpCK: expCK}
	case typeexp.TensorRef, typeexp.TensorRec:
		return typeexp.ExpRefDS{K: typeexp.TensorKind, ExpCK: expCK}
	case typeexp.LolliRef, typeexp.LolliRec:
		return typeexp.ExpRefDS{K: typeexp.LolliKind, ExpCK: expCK}
	case typeexp.PlusRef, typeexp.PlusRec:
		return typeexp.ExpRefDS{K: typeexp.PlusKind, ExpCK: expCK}
	case typeexp.WithRef, typeexp.WithRec:
		return typeexp.ExpRefDS{K: typeexp.WithKind, ExpCK: expCK}
	default:
		panic(typeexp.ErrRefTypeUnexpected(ref))
	}
}

func dataToExpRef(dto typeexp.ExpRefDS) (typeexp.ExpRef, error) {
	expCK, err := contkey.ConvertFromInt(dto.ExpCK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case typeexp.OneKind:
		return typeexp.OneRef{ExpCK: expCK}, nil
	case typeexp.LinkKind:
		return typeexp.LinkRef{ExpCK: expCK}, nil
	case typeexp.TensorKind:
		return typeexp.TensorRef{ExpCK: expCK}, nil
	case typeexp.LolliKind:
		return typeexp.LolliRef{ExpCK: expCK}, nil
	case typeexp.PlusKind:
		return typeexp.PlusRef{ExpCK: expCK}, nil
	case typeexp.WithKind:
		return typeexp.WithRef{ExpCK: expCK}, nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func dataToExpRec(dto typeexp.ExpRecDS) (typeexp.ExpRec, error) {
	states := make(map[int64]typeexp.StateDS, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ExpCK] = dto
	}
	return statesToExpRec(states, states[dto.ExpCK])
}

func dataFromExpRec(rec typeexp.ExpRec) typeexp.ExpRecDS {
	if rec == nil {
		panic("can't be nil")
	}
	dto := &typeexp.ExpRecDS{
		ExpCK:  contkey.ConvertToInt(rec.Key()),
		States: nil,
	}
	statesFromExpRec(0, rec, dto)
	return *dto
}

func statesToExpRec(states map[int64]typeexp.StateDS, st typeexp.StateDS) (typeexp.ExpRec, error) {
	expCK, err := contkey.ConvertFromInt(st.ExpCK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case typeexp.OneKind:
		return typeexp.OneRec{ExpCK: expCK}, nil
	case typeexp.LinkKind:
		typeQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return typeexp.LinkRec{ExpCK: expCK, TypeQN: typeQN}, nil
	case typeexp.TensorKind:
		b, err := statesToExpRec(states, states[st.Spec.Tensor.ValExpCK])
		if err != nil {
			return nil, err
		}
		c, err := statesToExpRec(states, states[st.Spec.Tensor.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typeexp.TensorRec{ExpCK: expCK, Val: b, Cont: c}, nil
	case typeexp.LolliKind:
		y, err := statesToExpRec(states, states[st.Spec.Lolli.ValExpCK])
		if err != nil {
			return nil, err
		}
		z, err := statesToExpRec(states, states[st.Spec.Lolli.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typeexp.LolliRec{ExpCK: expCK, Val: y, Cont: z}, nil
	case typeexp.PlusKind:
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
	case typeexp.WithKind:
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

func statesFromExpRec(fromID int64, r typeexp.ExpRec, dto *typeexp.ExpRecDS) int64 {
	expCK := contkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case typeexp.OneRec:
		st := typeexp.StateDS{ExpCK: expCK, K: typeexp.OneKind, SupExpCK: fromID}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.LinkRec:
		st := typeexp.StateDS{
			ExpCK:    expCK,
			K:        typeexp.LinkKind,
			SupExpCK: fromID,
			Spec: typeexp.ExpSpecDS{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.TensorRec:
		val := statesFromExpRec(expCK, rec.Val, dto)
		cont := statesFromExpRec(expCK, rec.Cont, dto)
		st := typeexp.StateDS{
			ExpCK:    expCK,
			K:        typeexp.TensorKind,
			SupExpCK: fromID,
			Spec: typeexp.ExpSpecDS{
				Tensor: &typeexp.ProdDS{ValExpCK: val, ContExpCK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.LolliRec:
		val := statesFromExpRec(expCK, rec.Val, dto)
		cont := statesFromExpRec(expCK, rec.Cont, dto)
		st := typeexp.StateDS{
			ExpCK:    expCK,
			K:        typeexp.LolliKind,
			SupExpCK: fromID,
			Spec: typeexp.ExpSpecDS{
				Lolli: &typeexp.ProdDS{ValExpCK: val, ContExpCK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.PlusRec:
		var choices []typeexp.SumDS
		for label, choice := range rec.Choices {
			cont := statesFromExpRec(expCK, choice, dto)
			choices = append(choices, typeexp.SumDS{LabQN: uniqsym.ConvertToString(label), ContExpCK: cont})
		}
		st := typeexp.StateDS{
			ExpCK:    expCK,
			K:        typeexp.PlusKind,
			SupExpCK: fromID,
			Spec:     typeexp.ExpSpecDS{Plus: choices},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.WithRec:
		var choices []typeexp.SumDS
		for label, choice := range rec.Choices {
			cont := statesFromExpRec(expCK, choice, dto)
			choices = append(choices, typeexp.SumDS{LabQN: uniqsym.ConvertToString(label), ContExpCK: cont})
		}
		st := typeexp.StateDS{
			ExpCK:    expCK,
			K:        typeexp.WithKind,
			SupExpCK: fromID,
			Spec:     typeexp.ExpSpecDS{With: choices},
		}
		dto.States = append(dto.States, st)
		return expCK
	default:
		panic(typeexp.ErrRecTypeUnexpected(r))
	}
}

func errUnexpectedKind(k typeexp.ExpKindDS) error {
	return fmt.Errorf("unexpected kind %v", k)
}
