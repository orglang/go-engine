package xactexp

import (
	"fmt"

	"golang.org/x/exp/maps"

	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"

	"github.com/orglang/go-sdk/adt/xactexp"
)

func ConvertSpecToRec(s ExpSpec) (ExpRec, error) {
	if s == nil {
		return nil, nil
	}
	switch spec := s.(type) {
	case OneSpec:
		return OneRec{ExpVK: valkey.One}, nil
	case LinkSpec:
		expVK, err := spec.XactQN.Key()
		if err != nil {
			return nil, err
		}
		return LinkRec{ExpVK: expVK, XactQN: spec.XactQN}, nil
	case WithSpec:
		choices := make(map[uniqsym.ADT]ExpRec, len(spec.Choices))
		keys := make([]valkey.ADT, len(spec.Choices)*2)
		for lab, spec := range spec.Choices {
			cont, err := ConvertSpecToRec(spec)
			if err != nil {
				return nil, err
			}
			choices[lab] = cont
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
		return WithRec{ExpVK: expVK, Choices: choices}, nil
	case PlusSpec:
		choices := make(map[uniqsym.ADT]ExpRec, len(spec.Choices))
		keys := make([]valkey.ADT, len(spec.Choices)*2)
		for lab, spec := range spec.Choices {
			cont, err := ConvertSpecToRec(spec)
			if err != nil {
				return nil, err
			}
			choices[lab] = cont
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
		return PlusRec{ExpVK: expVK, Choices: choices}, nil
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
		return LinkSpec{XactQN: rec.XactQN}
	case WithRec:
		choices := make(map[uniqsym.ADT]ExpSpec, len(rec.Choices))
		for lab, rec := range rec.Choices {
			choices[lab] = ConvertRecToSpec(rec)
		}
		return WithSpec{Choices: choices}
	case PlusRec:
		choices := make(map[uniqsym.ADT]ExpSpec, len(rec.Choices))
		for lab, st := range rec.Choices {
			choices[lab] = ConvertRecToSpec(st)
		}
		return PlusSpec{Choices: choices}
	default:
		panic(ErrRecTypeUnexpected(rec))
	}
}

func MsgFromExpSpec(s ExpSpec) xactexp.ExpSpec {
	switch spec := s.(type) {
	case OneSpec:
		return xactexp.ExpSpec{K: xactexp.One}
	case LinkSpec:
		return xactexp.ExpSpec{
			K:    xactexp.Link,
			Link: &xactexp.LinkSpec{XactQN: uniqsym.ConvertToString(spec.XactQN)}}
	case WithSpec:
		choices := make([]xactexp.ChoiceSpec, len(spec.Choices))
		for i, p := range maps.Keys(spec.Choices) {
			choices[i] = xactexp.ChoiceSpec{
				ProcQN: uniqsym.ConvertToString(p),
				ContES: MsgFromExpSpec(spec.Choices[p]),
			}
		}
		return xactexp.ExpSpec{K: xactexp.With, With: &xactexp.LaborSpec{Choices: choices}}
	case PlusSpec:
		choices := make([]xactexp.ChoiceSpec, len(spec.Choices))
		for i, p := range maps.Keys(spec.Choices) {
			choices[i] = xactexp.ChoiceSpec{
				ProcQN: uniqsym.ConvertToString(p),
				ContES: MsgFromExpSpec(spec.Choices[p]),
			}
		}
		return xactexp.ExpSpec{K: xactexp.Plus, Plus: &xactexp.LaborSpec{Choices: choices}}
	default:
		panic(ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpec(dto xactexp.ExpSpec) (ExpSpec, error) {
	switch dto.K {
	case xactexp.One:
		return OneSpec{}, nil
	case xactexp.Link:
		xactQN, err := uniqsym.ConvertFromString(dto.Link.XactQN)
		if err != nil {
			return nil, err
		}
		return LinkSpec{XactQN: xactQN}, nil
	case xactexp.Plus:
		choices := make(map[uniqsym.ADT]ExpSpec, len(dto.Plus.Choices))
		for _, ch := range dto.Plus.Choices {
			choice, err := MsgToExpSpec(ch.ContES)
			if err != nil {
				return nil, err
			}
			procQN, err := uniqsym.ConvertFromString(ch.ProcQN)
			if err != nil {
				return nil, err
			}
			choices[procQN] = choice
		}
		return PlusSpec{Choices: choices}, nil
	case xactexp.With:
		choices := make(map[uniqsym.ADT]ExpSpec, len(dto.With.Choices))
		for _, ch := range dto.With.Choices {
			choice, err := MsgToExpSpec(ch.ContES)
			if err != nil {
				return nil, err
			}
			procQN, err := uniqsym.ConvertFromString(ch.ProcQN)
			if err != nil {
				return nil, err
			}
			choices[procQN] = choice
		}
		return WithSpec{Choices: choices}, nil
	default:
		panic(xactexp.ErrKindUnexpected(dto.K))
	}
}

func msgFromExpRef(ref ExpRef) xactexp.ExpRef {
	expVK := valkey.ConvertToInteger(ref.Key())
	switch ref.(type) {
	case OneRef, OneRec:
		return xactexp.ExpRef{K: xactexp.One, ExpVK: expVK}
	case LinkRef, LinkRec:
		return xactexp.ExpRef{K: xactexp.Link, ExpVK: expVK}
	case PlusRef, PlusRec:
		return xactexp.ExpRef{K: xactexp.Plus, ExpVK: expVK}
	case WithRef, WithRec:
		return xactexp.ExpRef{K: xactexp.With, ExpVK: expVK}
	default:
		panic(ErrRefTypeUnexpected(ref))
	}
}

func msgToExpRef(dto xactexp.ExpRef) (ExpRef, error) {
	expVK, err := valkey.ConvertFromInteger(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case xactexp.One:
		return OneRef{expVK}, nil
	case xactexp.Link:
		return LinkRef{expVK}, nil
	case xactexp.Plus:
		return PlusRef{expVK}, nil
	case xactexp.With:
		return WithRef{expVK}, nil
	default:
		panic(xactexp.ErrKindUnexpected(dto.K))
	}
}

func dataFromExpRef(ref ExpRef) expRefDS {
	if ref == nil {
		panic("can't be nil")
	}
	expVK := valkey.ConvertToInteger(ref.Key())
	switch ref.(type) {
	case OneRef, OneRec:
		return expRefDS{K: oneExp, ExpVK: expVK}
	case LinkRef, LinkRec:
		return expRefDS{K: linkExp, ExpVK: expVK}
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
		xactQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return LinkRec{ExpVK: expVK, XactQN: xactQN}, nil
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
		st := stateDS{ExpVK: expVK, K: oneExp, SupExpVK: fromID}
		dto.States = append(dto.States, st)
		return expVK
	case LinkRec:
		st := stateDS{
			ExpVK:    expVK,
			K:        linkExp,
			SupExpVK: fromID,
			Spec: expSpecDS{
				Link: uniqsym.ConvertToString(rec.XactQN),
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
			SupExpVK: fromID,
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
			SupExpVK: fromID,
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
