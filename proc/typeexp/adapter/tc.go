package adapter

import (
	"fmt"

	"golang.org/x/exp/maps"

	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"

	sdktypeexp "github.com/orglang/go-sdk/proc/typeexp"

	typeexp "orglang/go-engine/proc/typeexp/core"
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
	expVK := valkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typeexp.OneRef, typeexp.OneRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.One, ExpVK: expVK}
	case typeexp.LinkRef, typeexp.LinkRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.Link, ExpVK: expVK}
	case typeexp.TensorRef, typeexp.TensorRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.Tensor, ExpVK: expVK}
	case typeexp.LolliRef, typeexp.LolliRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.Lolli, ExpVK: expVK}
	case typeexp.PlusRef, typeexp.PlusRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.Plus, ExpVK: expVK}
	case typeexp.WithRef, typeexp.WithRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.With, ExpVK: expVK}
	default:
		panic(typeexp.ErrRefTypeUnexpected(ref))
	}
}

func MsgToExpRef(dto sdktypeexp.ExpRef) (typeexp.ExpRef, error) {
	expVK, err := valkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case sdktypeexp.One:
		return typeexp.OneRef{ExpVK: expVK}, nil
	case sdktypeexp.Link:
		return typeexp.LinkRef{ExpVK: expVK}, nil
	case sdktypeexp.Tensor:
		return typeexp.TensorRef{ExpVK: expVK}, nil
	case sdktypeexp.Lolli:
		return typeexp.LolliRef{ExpVK: expVK}, nil
	case sdktypeexp.Plus:
		return typeexp.PlusRef{ExpVK: expVK}, nil
	case sdktypeexp.With:
		return typeexp.WithRef{ExpVK: expVK}, nil
	default:
		panic(sdktypeexp.ErrKindUnexpected(dto.K))
	}
}

func dataFromExpRef(ref typeexp.ExpRef) typeexp.ExpRefDS {
	if ref == nil {
		panic("can't be null")
	}
	expVK := valkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typeexp.OneRef, typeexp.OneRec:
		return typeexp.ExpRefDS{K: typeexp.OneKind, ExpVK: expVK}
	case typeexp.LinkRef, typeexp.LinkRec:
		return typeexp.ExpRefDS{K: typeexp.LinkKind, ExpVK: expVK}
	case typeexp.TensorRef, typeexp.TensorRec:
		return typeexp.ExpRefDS{K: typeexp.TensorKind, ExpVK: expVK}
	case typeexp.LolliRef, typeexp.LolliRec:
		return typeexp.ExpRefDS{K: typeexp.LolliKind, ExpVK: expVK}
	case typeexp.PlusRef, typeexp.PlusRec:
		return typeexp.ExpRefDS{K: typeexp.PlusKind, ExpVK: expVK}
	case typeexp.WithRef, typeexp.WithRec:
		return typeexp.ExpRefDS{K: typeexp.WithKind, ExpVK: expVK}
	default:
		panic(typeexp.ErrRefTypeUnexpected(ref))
	}
}

func dataToExpRef(dto typeexp.ExpRefDS) (typeexp.ExpRef, error) {
	expVK, err := valkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case typeexp.OneKind:
		return typeexp.OneRef{ExpVK: expVK}, nil
	case typeexp.LinkKind:
		return typeexp.LinkRef{ExpVK: expVK}, nil
	case typeexp.TensorKind:
		return typeexp.TensorRef{ExpVK: expVK}, nil
	case typeexp.LolliKind:
		return typeexp.LolliRef{ExpVK: expVK}, nil
	case typeexp.PlusKind:
		return typeexp.PlusRef{ExpVK: expVK}, nil
	case typeexp.WithKind:
		return typeexp.WithRef{ExpVK: expVK}, nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func dataToExpRec(dto typeexp.ExpRecDS) (typeexp.ExpRec, error) {
	states := make(map[int64]typeexp.StateDS, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ExpVK] = dto
	}
	return statesToExpRec(states, states[dto.ExpVK])
}

func dataFromExpRec(rec typeexp.ExpRec) typeexp.ExpRecDS {
	if rec == nil {
		panic("can't be nil")
	}
	dto := &typeexp.ExpRecDS{
		ExpVK:  valkey.ConvertToInt(rec.Key()),
		States: nil,
	}
	statesFromExpRec(0, rec, dto)
	return *dto
}

func statesToExpRec(states map[int64]typeexp.StateDS, st typeexp.StateDS) (typeexp.ExpRec, error) {
	expVK, err := valkey.ConvertFromInt(st.ExpVK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case typeexp.OneKind:
		return typeexp.OneRec{ExpVK: expVK}, nil
	case typeexp.LinkKind:
		typeQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return typeexp.LinkRec{ExpVK: expVK, TypeQN: typeQN}, nil
	case typeexp.TensorKind:
		b, err := statesToExpRec(states, states[st.Spec.Tensor.ValExpVK])
		if err != nil {
			return nil, err
		}
		c, err := statesToExpRec(states, states[st.Spec.Tensor.ContExpVK])
		if err != nil {
			return nil, err
		}
		return typeexp.TensorRec{ExpVK: expVK, Val: b, Cont: c}, nil
	case typeexp.LolliKind:
		y, err := statesToExpRec(states, states[st.Spec.Lolli.ValExpVK])
		if err != nil {
			return nil, err
		}
		z, err := statesToExpRec(states, states[st.Spec.Lolli.ContExpVK])
		if err != nil {
			return nil, err
		}
		return typeexp.LolliRec{ExpVK: expVK, Val: y, Cont: z}, nil
	case typeexp.PlusKind:
		choices := make(map[uniqsym.ADT]typeexp.ExpRec, len(st.Spec.Plus))
		for _, ch := range st.Spec.Plus {
			choice, err := statesToExpRec(states, states[ch.ContExpVK])
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(ch.LabQN)
			if err != nil {
				return nil, err
			}
			choices[label] = choice
		}
		return typeexp.PlusRec{ExpVK: expVK, Choices: choices}, nil
	case typeexp.WithKind:
		choices := make(map[uniqsym.ADT]typeexp.ExpRec, len(st.Spec.With))
		for _, ch := range st.Spec.With {
			choice, err := statesToExpRec(states, states[ch.ContExpVK])
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(ch.LabQN)
			if err != nil {
				return nil, err
			}
			choices[label] = choice
		}
		return typeexp.WithRec{ExpVK: expVK, Choices: choices}, nil
	default:
		panic(errUnexpectedKind(st.K))
	}
}

func statesFromExpRec(fromID int64, r typeexp.ExpRec, dto *typeexp.ExpRecDS) int64 {
	expVK := valkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case typeexp.OneRec:
		st := typeexp.StateDS{ExpVK: expVK, K: typeexp.OneKind, SupExpVK: fromID}
		dto.States = append(dto.States, st)
		return expVK
	case typeexp.LinkRec:
		st := typeexp.StateDS{
			ExpVK:    expVK,
			K:        typeexp.LinkKind,
			SupExpVK: fromID,
			Spec: typeexp.ExpSpecDS{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case typeexp.TensorRec:
		val := statesFromExpRec(expVK, rec.Val, dto)
		cont := statesFromExpRec(expVK, rec.Cont, dto)
		st := typeexp.StateDS{
			ExpVK:    expVK,
			K:        typeexp.TensorKind,
			SupExpVK: fromID,
			Spec: typeexp.ExpSpecDS{
				Tensor: &typeexp.ProdDS{ValExpVK: val, ContExpVK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case typeexp.LolliRec:
		val := statesFromExpRec(expVK, rec.Val, dto)
		cont := statesFromExpRec(expVK, rec.Cont, dto)
		st := typeexp.StateDS{
			ExpVK:    expVK,
			K:        typeexp.LolliKind,
			SupExpVK: fromID,
			Spec: typeexp.ExpSpecDS{
				Lolli: &typeexp.ProdDS{ValExpVK: val, ContExpVK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case typeexp.PlusRec:
		var choices []typeexp.SumDS
		for label, choice := range rec.Choices {
			cont := statesFromExpRec(expVK, choice, dto)
			choices = append(choices, typeexp.SumDS{LabQN: uniqsym.ConvertToString(label), ContExpVK: cont})
		}
		st := typeexp.StateDS{
			ExpVK:    expVK,
			K:        typeexp.PlusKind,
			SupExpVK: fromID,
			Spec:     typeexp.ExpSpecDS{Plus: choices},
		}
		dto.States = append(dto.States, st)
		return expVK
	case typeexp.WithRec:
		var choices []typeexp.SumDS
		for label, choice := range rec.Choices {
			cont := statesFromExpRec(expVK, choice, dto)
			choices = append(choices, typeexp.SumDS{LabQN: uniqsym.ConvertToString(label), ContExpVK: cont})
		}
		st := typeexp.StateDS{
			ExpVK:    expVK,
			K:        typeexp.WithKind,
			SupExpVK: fromID,
			Spec:     typeexp.ExpSpecDS{With: choices},
		}
		dto.States = append(dto.States, st)
		return expVK
	default:
		panic(typeexp.ErrRecTypeUnexpected(r))
	}
}

func errUnexpectedKind(k typeexp.ExpKindDS) error {
	return fmt.Errorf("unexpected kind %v", k)
}
