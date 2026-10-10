package adapter

import (
	"fmt"

	"golang.org/x/exp/maps"

	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"

	"github.com/orglang/go-sdk/proc/typeexp"

	proctypeexp "orglang/go-engine/proc/typeexp/core"
)

func MsgFromExpSpec(s proctypeexp.ExpSpec) typeexp.ExpSpec {
	switch spec := s.(type) {
	case proctypeexp.OneSpec:
		return typeexp.ExpSpec{K: typeexp.One}
	case proctypeexp.LinkSpec:
		return typeexp.ExpSpec{
			K:    typeexp.Link,
			Link: &typeexp.LinkSpec{TypeQN: uniqsym.ConvertToString(spec.TypeQN)}}
	case proctypeexp.TensorSpec:
		return typeexp.ExpSpec{
			K: typeexp.Tensor,
			Tensor: &typeexp.ProdSpec{
				Val:  MsgFromExpSpec(spec.Val),
				Cont: MsgFromExpSpec(spec.Cont),
			},
		}
	case proctypeexp.LolliSpec:
		return typeexp.ExpSpec{
			K: typeexp.Lolli,
			Lolli: &typeexp.ProdSpec{
				Val:  MsgFromExpSpec(spec.Val),
				Cont: MsgFromExpSpec(spec.Cont),
			},
		}
	case proctypeexp.WithSpec:
		choices := make([]typeexp.ChoiceSpec, len(spec.Choices))
		for i, l := range maps.Keys(spec.Choices) {
			choices[i] = typeexp.ChoiceSpec{
				LabQN: uniqsym.ConvertToString(l),
				Cont:  MsgFromExpSpec(spec.Choices[l]),
			}
		}
		return typeexp.ExpSpec{K: typeexp.With, With: &typeexp.SumSpec{Choices: choices}}
	case proctypeexp.PlusSpec:
		choices := make([]typeexp.ChoiceSpec, len(spec.Choices))
		for i, l := range maps.Keys(spec.Choices) {
			choices[i] = typeexp.ChoiceSpec{
				LabQN: uniqsym.ConvertToString(l),
				Cont:  MsgFromExpSpec(spec.Choices[l]),
			}
		}
		return typeexp.ExpSpec{K: typeexp.Plus, Plus: &typeexp.SumSpec{Choices: choices}}
	default:
		panic(proctypeexp.ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpec(dto typeexp.ExpSpec) (proctypeexp.ExpSpec, error) {
	switch dto.K {
	case typeexp.One:
		return proctypeexp.OneSpec{}, nil
	case typeexp.Link:
		typeQN, err := uniqsym.ConvertFromString(dto.Link.TypeQN)
		if err != nil {
			return nil, err
		}
		return proctypeexp.LinkSpec{TypeQN: typeQN}, nil
	case typeexp.Tensor:
		valES, err := MsgToExpSpec(dto.Tensor.Val)
		if err != nil {
			return nil, err
		}
		contES, err := MsgToExpSpec(dto.Tensor.Cont)
		if err != nil {
			return nil, err
		}
		return proctypeexp.TensorSpec{Val: valES, Cont: contES}, nil
	case typeexp.Lolli:
		valES, err := MsgToExpSpec(dto.Lolli.Val)
		if err != nil {
			return nil, err
		}
		contES, err := MsgToExpSpec(dto.Lolli.Cont)
		if err != nil {
			return nil, err
		}
		return proctypeexp.LolliSpec{Val: valES, Cont: contES}, nil
	case typeexp.Plus:
		choices := make(map[uniqsym.ADT]proctypeexp.ExpSpec, len(dto.Plus.Choices))
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
		return proctypeexp.PlusSpec{Choices: choices}, nil
	case typeexp.With:
		choices := make(map[uniqsym.ADT]proctypeexp.ExpSpec, len(dto.With.Choices))
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
		return proctypeexp.WithSpec{Choices: choices}, nil
	default:
		panic(typeexp.ErrKindUnexpected(dto.K))
	}
}

func MsgFromExpRef(ref proctypeexp.ExpRef) typeexp.ExpRef {
	expVK := valkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case proctypeexp.OneRef, proctypeexp.OneRec:
		return typeexp.ExpRef{K: typeexp.One, ExpVK: expVK}
	case proctypeexp.LinkRef, proctypeexp.LinkRec:
		return typeexp.ExpRef{K: typeexp.Link, ExpVK: expVK}
	case proctypeexp.TensorRef, proctypeexp.TensorRec:
		return typeexp.ExpRef{K: typeexp.Tensor, ExpVK: expVK}
	case proctypeexp.LolliRef, proctypeexp.LolliRec:
		return typeexp.ExpRef{K: typeexp.Lolli, ExpVK: expVK}
	case proctypeexp.PlusRef, proctypeexp.PlusRec:
		return typeexp.ExpRef{K: typeexp.Plus, ExpVK: expVK}
	case proctypeexp.WithRef, proctypeexp.WithRec:
		return typeexp.ExpRef{K: typeexp.With, ExpVK: expVK}
	default:
		panic(proctypeexp.ErrRefTypeUnexpected(ref))
	}
}

func MsgToExpRef(dto typeexp.ExpRef) (proctypeexp.ExpRef, error) {
	expVK, err := valkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case typeexp.One:
		return proctypeexp.OneRef{ExpVK: expVK}, nil
	case typeexp.Link:
		return proctypeexp.LinkRef{ExpVK: expVK}, nil
	case typeexp.Tensor:
		return proctypeexp.TensorRef{ExpVK: expVK}, nil
	case typeexp.Lolli:
		return proctypeexp.LolliRef{ExpVK: expVK}, nil
	case typeexp.Plus:
		return proctypeexp.PlusRef{ExpVK: expVK}, nil
	case typeexp.With:
		return proctypeexp.WithRef{ExpVK: expVK}, nil
	default:
		panic(typeexp.ErrKindUnexpected(dto.K))
	}
}

func dataFromExpRef(ref proctypeexp.ExpRef) proctypeexp.ExpRefDS {
	if ref == nil {
		panic("can't be null")
	}
	expVK := valkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case proctypeexp.OneRef, proctypeexp.OneRec:
		return proctypeexp.ExpRefDS{K: proctypeexp.OneKind, ExpVK: expVK}
	case proctypeexp.LinkRef, proctypeexp.LinkRec:
		return proctypeexp.ExpRefDS{K: proctypeexp.LinkKind, ExpVK: expVK}
	case proctypeexp.TensorRef, proctypeexp.TensorRec:
		return proctypeexp.ExpRefDS{K: proctypeexp.TensorKind, ExpVK: expVK}
	case proctypeexp.LolliRef, proctypeexp.LolliRec:
		return proctypeexp.ExpRefDS{K: proctypeexp.LolliKind, ExpVK: expVK}
	case proctypeexp.PlusRef, proctypeexp.PlusRec:
		return proctypeexp.ExpRefDS{K: proctypeexp.PlusKind, ExpVK: expVK}
	case proctypeexp.WithRef, proctypeexp.WithRec:
		return proctypeexp.ExpRefDS{K: proctypeexp.WithKind, ExpVK: expVK}
	default:
		panic(proctypeexp.ErrRefTypeUnexpected(ref))
	}
}

func dataToExpRef(dto proctypeexp.ExpRefDS) (proctypeexp.ExpRef, error) {
	expVK, err := valkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case proctypeexp.OneKind:
		return proctypeexp.OneRef{ExpVK: expVK}, nil
	case proctypeexp.LinkKind:
		return proctypeexp.LinkRef{ExpVK: expVK}, nil
	case proctypeexp.TensorKind:
		return proctypeexp.TensorRef{ExpVK: expVK}, nil
	case proctypeexp.LolliKind:
		return proctypeexp.LolliRef{ExpVK: expVK}, nil
	case proctypeexp.PlusKind:
		return proctypeexp.PlusRef{ExpVK: expVK}, nil
	case proctypeexp.WithKind:
		return proctypeexp.WithRef{ExpVK: expVK}, nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func dataToExpRec(dto proctypeexp.ExpRecDS) (proctypeexp.ExpRec, error) {
	states := make(map[int64]proctypeexp.StateDS, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ExpVK] = dto
	}
	return statesToExpRec(states, states[dto.ExpVK])
}

func dataFromExpRec(rec proctypeexp.ExpRec) proctypeexp.ExpRecDS {
	if rec == nil {
		panic("can't be nil")
	}
	dto := &proctypeexp.ExpRecDS{
		ExpVK:  valkey.ConvertToInt(rec.Key()),
		States: nil,
	}
	statesFromExpRec(0, rec, dto)
	return *dto
}

func statesToExpRec(states map[int64]proctypeexp.StateDS, st proctypeexp.StateDS) (proctypeexp.ExpRec, error) {
	expVK, err := valkey.ConvertFromInt(st.ExpVK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case proctypeexp.OneKind:
		return proctypeexp.OneRec{ExpVK: expVK}, nil
	case proctypeexp.LinkKind:
		typeQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return proctypeexp.LinkRec{ExpVK: expVK, TypeQN: typeQN}, nil
	case proctypeexp.TensorKind:
		b, err := statesToExpRec(states, states[st.Spec.Tensor.ValExpVK])
		if err != nil {
			return nil, err
		}
		c, err := statesToExpRec(states, states[st.Spec.Tensor.ContExpVK])
		if err != nil {
			return nil, err
		}
		return proctypeexp.TensorRec{ExpVK: expVK, Val: b, Cont: c}, nil
	case proctypeexp.LolliKind:
		y, err := statesToExpRec(states, states[st.Spec.Lolli.ValExpVK])
		if err != nil {
			return nil, err
		}
		z, err := statesToExpRec(states, states[st.Spec.Lolli.ContExpVK])
		if err != nil {
			return nil, err
		}
		return proctypeexp.LolliRec{ExpVK: expVK, Val: y, Cont: z}, nil
	case proctypeexp.PlusKind:
		choices := make(map[uniqsym.ADT]proctypeexp.ExpRec, len(st.Spec.Plus))
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
		return proctypeexp.PlusRec{ExpVK: expVK, Choices: choices}, nil
	case proctypeexp.WithKind:
		choices := make(map[uniqsym.ADT]proctypeexp.ExpRec, len(st.Spec.With))
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
		return proctypeexp.WithRec{ExpVK: expVK, Choices: choices}, nil
	default:
		panic(errUnexpectedKind(st.K))
	}
}

func statesFromExpRec(fromID int64, r proctypeexp.ExpRec, dto *proctypeexp.ExpRecDS) int64 {
	expVK := valkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case proctypeexp.OneRec:
		st := proctypeexp.StateDS{ExpVK: expVK, K: proctypeexp.OneKind, SupExpVK: fromID}
		dto.States = append(dto.States, st)
		return expVK
	case proctypeexp.LinkRec:
		st := proctypeexp.StateDS{
			ExpVK:    expVK,
			K:        proctypeexp.LinkKind,
			SupExpVK: fromID,
			Spec: proctypeexp.ExpSpecDS{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case proctypeexp.TensorRec:
		val := statesFromExpRec(expVK, rec.Val, dto)
		cont := statesFromExpRec(expVK, rec.Cont, dto)
		st := proctypeexp.StateDS{
			ExpVK:    expVK,
			K:        proctypeexp.TensorKind,
			SupExpVK: fromID,
			Spec: proctypeexp.ExpSpecDS{
				Tensor: &proctypeexp.ProdDS{ValExpVK: val, ContExpVK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case proctypeexp.LolliRec:
		val := statesFromExpRec(expVK, rec.Val, dto)
		cont := statesFromExpRec(expVK, rec.Cont, dto)
		st := proctypeexp.StateDS{
			ExpVK:    expVK,
			K:        proctypeexp.LolliKind,
			SupExpVK: fromID,
			Spec: proctypeexp.ExpSpecDS{
				Lolli: &proctypeexp.ProdDS{ValExpVK: val, ContExpVK: cont},
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case proctypeexp.PlusRec:
		var choices []proctypeexp.SumDS
		for label, choice := range rec.Choices {
			cont := statesFromExpRec(expVK, choice, dto)
			choices = append(choices, proctypeexp.SumDS{LabQN: uniqsym.ConvertToString(label), ContExpVK: cont})
		}
		st := proctypeexp.StateDS{
			ExpVK:    expVK,
			K:        proctypeexp.PlusKind,
			SupExpVK: fromID,
			Spec:     proctypeexp.ExpSpecDS{Plus: choices},
		}
		dto.States = append(dto.States, st)
		return expVK
	case proctypeexp.WithRec:
		var choices []proctypeexp.SumDS
		for label, choice := range rec.Choices {
			cont := statesFromExpRec(expVK, choice, dto)
			choices = append(choices, proctypeexp.SumDS{LabQN: uniqsym.ConvertToString(label), ContExpVK: cont})
		}
		st := proctypeexp.StateDS{
			ExpVK:    expVK,
			K:        proctypeexp.WithKind,
			SupExpVK: fromID,
			Spec:     proctypeexp.ExpSpecDS{With: choices},
		}
		dto.States = append(dto.States, st)
		return expVK
	default:
		panic(proctypeexp.ErrRecTypeUnexpected(r))
	}
}

func errUnexpectedKind(k proctypeexp.ExpKindDS) error {
	return fmt.Errorf("unexpected kind %v", k)
}
