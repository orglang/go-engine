package xactexp

import (
	"database/sql"
	"fmt"

	"golang.org/x/exp/maps"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/xactexp"
)

func ConvertSpecToRec(s ExpSpec) ExpRec {
	if s == nil {
		return nil
	}
	switch spec := s.(type) {
	case OneSpec:
		return OneRec{ExpID: identity.New()}
	case LinkSpec:
		return LinkRec{ExpID: identity.New(), XactQN: spec.XactQN}
	case WithSpec:
		choices := make(map[uniqsym.ADT]ExpRec, len(spec.Choices))
		for lab, st := range spec.Choices {
			choices[lab] = ConvertSpecToRec(st)
		}
		return WithRec{ExpID: identity.New(), Zs: choices}
	case PlusSpec:
		choices := make(map[uniqsym.ADT]ExpRec, len(spec.Choices))
		for lab, rec := range spec.Choices {
			choices[lab] = ConvertSpecToRec(rec)
		}
		return PlusRec{ExpID: identity.New(), Zs: choices}
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
		choices := make(map[uniqsym.ADT]ExpSpec, len(rec.Zs))
		for lab, rec := range rec.Zs {
			choices[lab] = ConvertRecToSpec(rec)
		}
		return WithSpec{Choices: choices}
	case PlusRec:
		choices := make(map[uniqsym.ADT]ExpSpec, len(rec.Zs))
		for lab, st := range rec.Zs {
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

func msgFromExpRef(r ExpRef) xactexp.ExpRef {
	ident := r.Ident().String()
	switch r.(type) {
	case OneRef, OneRec:
		return xactexp.ExpRef{K: xactexp.One, ExpID: ident}
	case LinkRef, LinkRec:
		return xactexp.ExpRef{K: xactexp.Link, ExpID: ident}
	case PlusRef, PlusRec:
		return xactexp.ExpRef{K: xactexp.Plus, ExpID: ident}
	case WithRef, WithRec:
		return xactexp.ExpRef{K: xactexp.With, ExpID: ident}
	default:
		panic(ErrRefTypeUnexpected(r))
	}
}

func msgToExpRef(dto xactexp.ExpRef) (ExpRef, error) {
	expID, err := identity.ConvertFromString(dto.ExpID)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case xactexp.One:
		return OneRef{expID}, nil
	case xactexp.Link:
		return LinkRef{expID}, nil
	case xactexp.Plus:
		return PlusRef{expID}, nil
	case xactexp.With:
		return WithRef{expID}, nil
	default:
		panic(xactexp.ErrKindUnexpected(dto.K))
	}
}

func dataFromExpRef(ref ExpRef) expRefDS {
	if ref == nil {
		panic("can't be nil")
	}
	expID := ref.Ident().String()
	switch ref.(type) {
	case OneRef, OneRec:
		return expRefDS{K: oneExp, ExpID: expID}
	case LinkRef, LinkRec:
		return expRefDS{K: linkExp, ExpID: expID}
	case PlusRef, PlusRec:
		return expRefDS{K: plusExp, ExpID: expID}
	case WithRef, WithRec:
		return expRefDS{K: withExp, ExpID: expID}
	default:
		panic(ErrRefTypeUnexpected(ref))
	}
}

func dataToExpRef(dto expRefDS) (ExpRef, error) {
	expID, err := identity.ConvertFromString(dto.ExpID)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case oneExp:
		return OneRef{expID}, nil
	case linkExp:
		return LinkRef{expID}, nil
	case plusExp:
		return PlusRef{expID}, nil
	case withExp:
		return WithRef{expID}, nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func dataToExpRec(dto expRecDS) (ExpRec, error) {
	states := make(map[string]stateDS, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ExpID] = dto
	}
	return statesToExpRec(states, states[dto.ExpID])
}

func dataFromExpRec(rec ExpRec) expRecDS {
	if rec == nil {
		panic("can't be nil")
	}
	dto := &expRecDS{
		ExpID:  rec.Ident().String(),
		States: nil,
	}
	statesFromExpRec("", rec, dto)
	return *dto
}

func statesToExpRec(states map[string]stateDS, st stateDS) (ExpRec, error) {
	expID, err := identity.ConvertFromString(st.ExpID)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case oneExp:
		return OneRec{ExpID: expID}, nil
	case linkExp:
		xactQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return LinkRec{ExpID: expID, XactQN: xactQN}, nil
	case plusExp:
		choices := make(map[uniqsym.ADT]ExpRec, len(st.Spec.Plus))
		for _, ch := range st.Spec.Plus {
			choice, err := statesToExpRec(states, states[ch.ContES])
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(ch.LabQN)
			if err != nil {
				return nil, err
			}
			choices[label] = choice
		}
		return PlusRec{ExpID: expID, Zs: choices}, nil
	case withExp:
		choices := make(map[uniqsym.ADT]ExpRec, len(st.Spec.With))
		for _, ch := range st.Spec.With {
			choice, err := statesToExpRec(states, states[ch.ContES])
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(ch.LabQN)
			if err != nil {
				return nil, err
			}
			choices[label] = choice
		}
		return WithRec{ExpID: expID, Zs: choices}, nil
	default:
		panic(errUnexpectedKind(st.K))
	}
}

func statesFromExpRec(from string, r ExpRec, dto *expRecDS) string {
	var fromID sql.NullString
	if len(from) > 0 {
		fromID = sql.NullString{String: from, Valid: true}
	}
	expID := r.Ident().String()
	switch rec := r.(type) {
	case OneRec:
		st := stateDS{ExpID: expID, K: oneExp, SupExpID: fromID}
		dto.States = append(dto.States, st)
		return expID
	case LinkRec:
		st := stateDS{
			ExpID:    expID,
			K:        linkExp,
			SupExpID: fromID,
			Spec: expSpecDS{
				Link: uniqsym.ConvertToString(rec.XactQN),
			},
		}
		dto.States = append(dto.States, st)
		return expID
	case PlusRec:
		var choices []sumDS
		for label, choice := range rec.Zs {
			cont := statesFromExpRec(expID, choice, dto)
			choices = append(choices, sumDS{uniqsym.ConvertToString(label), cont})
		}
		st := stateDS{
			ExpID:    expID,
			K:        plusExp,
			SupExpID: fromID,
			Spec:     expSpecDS{Plus: choices},
		}
		dto.States = append(dto.States, st)
		return expID
	case WithRec:
		var choices []sumDS
		for label, choice := range rec.Zs {
			cont := statesFromExpRec(expID, choice, dto)
			choices = append(choices, sumDS{uniqsym.ConvertToString(label), cont})
		}
		st := stateDS{
			ExpID:    expID,
			K:        withExp,
			SupExpID: fromID,
			Spec:     expSpecDS{With: choices},
		}
		dto.States = append(dto.States, st)
		return expID
	default:
		panic(ErrRecTypeUnexpected(r))
	}
}

func errUnexpectedKind(k expKindDS) error {
	return fmt.Errorf("unexpected kind %q", k)
}
