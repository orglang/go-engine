package typeexp

import (
	"fmt"

	"golang.org/x/exp/maps"

	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"

	"github.com/orglang/go-sdk/adt/typeexp"
)

func ConvertSpecToRec(s ExpSpec) (ExpRec, error) {
	if s == nil {
		return nil, nil
	}
	switch spec := s.(type) {
	case OneSpec:
		return OneRec{ExpVK: valkey.One}, nil
	case LinkSpec:
		expVK, err := spec.TypeQN.Key()
		if err != nil {
			return nil, err
		}
		return LinkRec{ExpVK: expVK, TypeQN: spec.TypeQN}, nil
	case TensorSpec:
		val, err := ConvertSpecToRec(spec.Val)
		if err != nil {
			return nil, err
		}
		cont, err := ConvertSpecToRec(spec.Cont)
		if err != nil {
			return nil, err
		}
		key := val.Key() + cont.Key()
		return TensorRec{ExpVK: key, Val: val, Cont: cont}, nil
	case LolliSpec:
		val, err := ConvertSpecToRec(spec.Val)
		if err != nil {
			return nil, err
		}
		cont, err := ConvertSpecToRec(spec.Cont)
		if err != nil {
			return nil, err
		}
		key := val.Key() + cont.Key()
		return LolliRec{ExpVK: key, Val: val, Cont: cont}, nil
	case WithSpec:
		conts := make(map[uniqsym.ADT]ExpRec, len(spec.Choices))
		keys := make([]valkey.ADT, len(spec.Choices)*2)
		for lab, spec := range spec.Choices {
			cont, err := ConvertSpecToRec(spec)
			if err != nil {
				return nil, err
			}
			conts[lab] = cont
			labVK, err := lab.Key()
			if err != nil {
				return nil, err
			}
			keys = append(keys, labVK, cont.Key())
		}
		expVK, err := valkey.Compose(keys...)
		if err != nil {
			return nil, err
		}
		return WithRec{ExpVK: expVK, Choices: conts}, nil
	case PlusSpec:
		conts := make(map[uniqsym.ADT]ExpRec, len(spec.Choices))
		keys := make([]valkey.ADT, len(spec.Choices)*2)
		for lab, spec := range spec.Choices {
			cont, err := ConvertSpecToRec(spec)
			if err != nil {
				return nil, err
			}
			conts[lab] = cont
			labVK, err := lab.Key()
			if err != nil {
				return nil, err
			}
			keys = append(keys, labVK, cont.Key())
		}
		expVK, err := valkey.Compose(keys...)
		if err != nil {
			return nil, err
		}
		return PlusRec{ExpVK: expVK, Choices: conts}, nil
	default:
		panic(ErrSpecTypeUnexpected(spec))
	}
}

func ConvertRecToSpec(r ExpRec) ExpSpec {
	if r == nil {
		return nil
	}
	switch rec := r.(type) {
	case OneRec:
		return OneSpec{}
	case LinkRec:
		return LinkSpec{TypeQN: rec.TypeQN}
	case TensorRec:
		return TensorSpec{
			Val:  ConvertRecToSpec(rec.Val),
			Cont: ConvertRecToSpec(rec.Cont),
		}
	case LolliRec:
		return LolliSpec{
			Val:  ConvertRecToSpec(rec.Val),
			Cont: ConvertRecToSpec(rec.Cont),
		}
	case WithRec:
		choices := make(map[uniqsym.ADT]ExpSpec, len(rec.Choices))
		for lab, cont := range rec.Choices {
			choices[lab] = ConvertRecToSpec(cont)
		}
		return WithSpec{Choices: choices}
	case PlusRec:
		choices := make(map[uniqsym.ADT]ExpSpec, len(rec.Choices))
		for lab, cont := range rec.Choices {
			choices[lab] = ConvertRecToSpec(cont)
		}
		return PlusSpec{Choices: choices}
	default:
		panic(ErrRecTypeUnexpected(rec))
	}
}

func MsgFromExpSpec(s ExpSpec) typeexp.ExpSpec {
	switch spec := s.(type) {
	case OneSpec:
		return typeexp.ExpSpec{K: typeexp.One}
	case LinkSpec:
		return typeexp.ExpSpec{
			K:    typeexp.Link,
			Link: &typeexp.LinkSpec{TypeQN: uniqsym.ConvertToString(spec.TypeQN)}}
	case TensorSpec:
		return typeexp.ExpSpec{
			K: typeexp.Tensor,
			Tensor: &typeexp.ProdSpec{
				Val:  MsgFromExpSpec(spec.Val),
				Cont: MsgFromExpSpec(spec.Cont),
			},
		}
	case LolliSpec:
		return typeexp.ExpSpec{
			K: typeexp.Lolli,
			Lolli: &typeexp.ProdSpec{
				Val:  MsgFromExpSpec(spec.Val),
				Cont: MsgFromExpSpec(spec.Cont),
			},
		}
	case WithSpec:
		choices := make([]typeexp.ChoiceSpec, len(spec.Choices))
		for i, l := range maps.Keys(spec.Choices) {
			choices[i] = typeexp.ChoiceSpec{
				LabQN: uniqsym.ConvertToString(l),
				Cont:  MsgFromExpSpec(spec.Choices[l]),
			}
		}
		return typeexp.ExpSpec{K: typeexp.With, With: &typeexp.SumSpec{Choices: choices}}
	case PlusSpec:
		choices := make([]typeexp.ChoiceSpec, len(spec.Choices))
		for i, l := range maps.Keys(spec.Choices) {
			choices[i] = typeexp.ChoiceSpec{
				LabQN: uniqsym.ConvertToString(l),
				Cont:  MsgFromExpSpec(spec.Choices[l]),
			}
		}
		return typeexp.ExpSpec{K: typeexp.Plus, Plus: &typeexp.SumSpec{Choices: choices}}
	default:
		panic(ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpec(dto typeexp.ExpSpec) (ExpSpec, error) {
	switch dto.K {
	case typeexp.One:
		return OneSpec{}, nil
	case typeexp.Link:
		typeQN, err := uniqsym.ConvertFromString(dto.Link.TypeQN)
		if err != nil {
			return nil, err
		}
		return LinkSpec{TypeQN: typeQN}, nil
	case typeexp.Tensor:
		valES, err := MsgToExpSpec(dto.Tensor.Val)
		if err != nil {
			return nil, err
		}
		contES, err := MsgToExpSpec(dto.Tensor.Cont)
		if err != nil {
			return nil, err
		}
		return TensorSpec{Val: valES, Cont: contES}, nil
	case typeexp.Lolli:
		valES, err := MsgToExpSpec(dto.Lolli.Val)
		if err != nil {
			return nil, err
		}
		contES, err := MsgToExpSpec(dto.Lolli.Cont)
		if err != nil {
			return nil, err
		}
		return LolliSpec{Val: valES, Cont: contES}, nil
	case typeexp.Plus:
		choices := make(map[uniqsym.ADT]ExpSpec, len(dto.Plus.Choices))
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
		return PlusSpec{Choices: choices}, nil
	case typeexp.With:
		choices := make(map[uniqsym.ADT]ExpSpec, len(dto.With.Choices))
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
		return WithSpec{Choices: choices}, nil
	default:
		panic(typeexp.ErrKindUnexpected(dto.K))
	}
}

func MsgFromExpRef(ref ExpRef) typeexp.ExpRef {
	expVK := valkey.ConvertToInteger(ref.Key())
	switch ref.(type) {
	case OneRef, OneRec:
		return typeexp.ExpRef{K: typeexp.One, ExpVK: expVK}
	case LinkRef, LinkRec:
		return typeexp.ExpRef{K: typeexp.Link, ExpVK: expVK}
	case TensorRef, TensorRec:
		return typeexp.ExpRef{K: typeexp.Tensor, ExpVK: expVK}
	case LolliRef, LolliRec:
		return typeexp.ExpRef{K: typeexp.Lolli, ExpVK: expVK}
	case PlusRef, PlusRec:
		return typeexp.ExpRef{K: typeexp.Plus, ExpVK: expVK}
	case WithRef, WithRec:
		return typeexp.ExpRef{K: typeexp.With, ExpVK: expVK}
	default:
		panic(ErrRefTypeUnexpected(ref))
	}
}

func MsgToExpRef(dto typeexp.ExpRef) (ExpRef, error) {
	expVK, err := valkey.ConvertFromInteger(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case typeexp.One:
		return OneRef{expVK}, nil
	case typeexp.Link:
		return LinkRef{expVK}, nil
	case typeexp.Tensor:
		return TensorRef{expVK}, nil
	case typeexp.Lolli:
		return LolliRef{expVK}, nil
	case typeexp.Plus:
		return PlusRef{expVK}, nil
	case typeexp.With:
		return WithRef{expVK}, nil
	default:
		panic(typeexp.ErrKindUnexpected(dto.K))
	}
}

func dataFromExpRef(ref ExpRef) expRefDS {
	if ref == nil {
		panic("can't be null")
	}
	expVK := valkey.ConvertToInteger(ref.Key())
	switch ref.(type) {
	case OneRef, OneRec:
		return expRefDS{K: oneExp, ExpVK: expVK}
	case LinkRef, LinkRec:
		return expRefDS{K: linkExp, ExpVK: expVK}
	case TensorRef, TensorRec:
		return expRefDS{K: tensorExp, ExpVK: expVK}
	case LolliRef, LolliRec:
		return expRefDS{K: lolliExp, ExpVK: expVK}
	case PlusRef, PlusRec:
		return expRefDS{K: plusExp, ExpVK: expVK}
	case WithRef, WithRec:
		return expRefDS{K: withExp, ExpVK: expVK}
	default:
		panic(ErrRefTypeUnexpected(ref))
	}
}

func dataToExpRef(dto expRefDS) (ExpRef, error) {
	expVK, err := valkey.ConvertFromInteger(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case oneExp:
		return OneRef{expVK}, nil
	case linkExp:
		return LinkRef{expVK}, nil
	case tensorExp:
		return TensorRef{expVK}, nil
	case lolliExp:
		return LolliRef{expVK}, nil
	case plusExp:
		return PlusRef{expVK}, nil
	case withExp:
		return WithRef{expVK}, nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func dataToExpRec(dto expRecDS) (ExpRec, error) {
	states := make(map[int64]stateDS, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ExpVK] = dto
	}
	return statesToExpRec(states, states[dto.ExpVK])
}

func dataFromExpRec(rec ExpRec) expRecDS {
	if rec == nil {
		panic("can't be nil")
	}
	dto := &expRecDS{
		ExpVK:  valkey.ConvertToInteger(rec.Key()),
		States: nil,
	}
	statesFromExpRec(0, rec, dto)
	return *dto
}

func statesToExpRec(states map[int64]stateDS, st stateDS) (ExpRec, error) {
	expVK, err := valkey.ConvertFromInteger(st.ExpVK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case oneExp:
		return OneRec{ExpVK: expVK}, nil
	case linkExp:
		typeQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return LinkRec{ExpVK: expVK, TypeQN: typeQN}, nil
	case tensorExp:
		b, err := statesToExpRec(states, states[st.Spec.Tensor.ValExpVK])
		if err != nil {
			return nil, err
		}
		c, err := statesToExpRec(states, states[st.Spec.Tensor.ContExpVK])
		if err != nil {
			return nil, err
		}
		return TensorRec{ExpVK: expVK, Val: b, Cont: c}, nil
	case lolliExp:
		y, err := statesToExpRec(states, states[st.Spec.Lolli.ValExpVK])
		if err != nil {
			return nil, err
		}
		z, err := statesToExpRec(states, states[st.Spec.Lolli.ContExpVK])
		if err != nil {
			return nil, err
		}
		return LolliRec{ExpVK: expVK, Val: y, Cont: z}, nil
	case plusExp:
		choices := make(map[uniqsym.ADT]ExpRec, len(st.Spec.Plus))
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
		return PlusRec{ExpVK: expVK, Choices: choices}, nil
	case withExp:
		choices := make(map[uniqsym.ADT]ExpRec, len(st.Spec.With))
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
		return WithRec{ExpVK: expVK, Choices: choices}, nil
	default:
		panic(errUnexpectedKind(st.K))
	}
}

func statesFromExpRec(fromID int64, r ExpRec, dto *expRecDS) int64 {
	expVK := valkey.ConvertToInteger(r.Key())
	switch rec := r.(type) {
	case OneRec:
		st := stateDS{ExpVK: expVK, K: oneExp, SupExpSK: fromID}
		dto.States = append(dto.States, st)
		return expVK
	case LinkRec:
		st := stateDS{
			ExpVK:    expVK,
			K:        linkExp,
			SupExpSK: fromID,
			Spec: expSpecDS{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case TensorRec:
		val := statesFromExpRec(expVK, rec.Val, dto)
		cont := statesFromExpRec(expVK, rec.Cont, dto)
		st := stateDS{
			ExpVK:    expVK,
			K:        tensorExp,
			SupExpSK: fromID,
			Spec: expSpecDS{
				Tensor: &prodDS{val, cont},
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case LolliRec:
		val := statesFromExpRec(expVK, rec.Val, dto)
		cont := statesFromExpRec(expVK, rec.Cont, dto)
		st := stateDS{
			ExpVK:    expVK,
			K:        lolliExp,
			SupExpSK: fromID,
			Spec: expSpecDS{
				Lolli: &prodDS{val, cont},
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case PlusRec:
		var choices []sumDS
		for label, choice := range rec.Choices {
			cont := statesFromExpRec(expVK, choice, dto)
			choices = append(choices, sumDS{uniqsym.ConvertToString(label), cont})
		}
		st := stateDS{
			ExpVK:    expVK,
			K:        plusExp,
			SupExpSK: fromID,
			Spec:     expSpecDS{Plus: choices},
		}
		dto.States = append(dto.States, st)
		return expVK
	case WithRec:
		var choices []sumDS
		for label, choice := range rec.Choices {
			cont := statesFromExpRec(expVK, choice, dto)
			choices = append(choices, sumDS{uniqsym.ConvertToString(label), cont})
		}
		st := stateDS{
			ExpVK:    expVK,
			K:        withExp,
			SupExpSK: fromID,
			Spec:     expSpecDS{With: choices},
		}
		dto.States = append(dto.States, st)
		return expVK
	default:
		panic(ErrRecTypeUnexpected(r))
	}
}

func errUnexpectedKind(k expKindDS) error {
	return fmt.Errorf("unexpected kind %q", k)
}
