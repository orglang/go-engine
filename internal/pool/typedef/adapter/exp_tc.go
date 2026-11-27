package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"

	typedef "orglang/go-engine/internal/pool/typedef/core"
)

func dataFromExpRef(ref typedef.ExpRef) ExpRef {
	if ref == nil {
		panic("can't be nil")
	}
	expCK := contkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typedef.OneRef, typedef.OneRec:
		return ExpRef{K: OneKind, ExpCK: expCK}
	case typedef.LinkRef, typedef.LinkRec:
		return ExpRef{K: LinkKind, ExpCK: expCK}
	case typedef.PlusRef, typedef.PlusRec:
		return ExpRef{K: PlusKind, ExpCK: expCK}
	case typedef.WithRef, typedef.WithRec:
		return ExpRef{K: WithKind, ExpCK: expCK}
	default:
		panic(typedef.ErrRefTypeUnexpected(ref))
	}
}

func dataToExpRef(dto ExpRef) (typedef.ExpRef, error) {
	expCK, err := contkey.ConvertFromInt(dto.ExpCK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case OneKind:
		return typedef.OneRef{ExpCK: expCK}, nil
	case LinkKind:
		return typedef.LinkRef{ExpCK: expCK}, nil
	case PlusKind:
		return typedef.PlusRef{ExpCK: expCK}, nil
	case WithKind:
		return typedef.WithRef{ExpCK: expCK}, nil
	default:
		panic(errExpKindUnexpected(dto.K))
	}
}

func dataToExpRec(dto ExpRec) (typedef.ExpRec, error) {
	states := make(map[int64]State, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ExpCK] = dto
	}
	return statesToExpRec(states, states[dto.ExpCK])
}

func dataFromExpRec(rec typedef.ExpRec) ExpRec {
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

func statesToExpRec(states map[int64]State, st State) (typedef.ExpRec, error) {
	expCK, err := contkey.ConvertFromInt(st.ExpCK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case OneKind:
		return typedef.OneRec{ExpCK: expCK}, nil
	case LinkKind:
		xactQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return typedef.LinkRec{ExpCK: expCK, TypeQN: xactQN}, nil
	case PlusKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.Plus.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := statesToExpRec(states, states[st.Spec.Plus.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.PlusRec{ExpCK: expCK, ProcQNs: procQNs, ContExp: contExp}, nil
	case WithKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.With.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := statesToExpRec(states, states[st.Spec.With.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.WithRec{ExpCK: expCK, ProcQNs: procQNs, ContExp: contExp}, nil
	case UpKind:
		contExp, err := statesToExpRec(states, states[st.Spec.Up.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.UpRec{ExpCK: expCK, ContExp: contExp}, nil
	case DownKind:
		contExp, err := statesToExpRec(states, states[st.Spec.Down.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typedef.DownRec{ExpCK: expCK, ContExp: contExp}, nil
	default:
		panic(errExpKindUnexpected(st.K))
	}
}

func statesFromExpRec(supExpCK int64, r typedef.ExpRec, dto *ExpRec) int64 {
	expCK := contkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case typedef.OneRec:
		st := State{ExpCK: expCK, K: OneKind, SupExpCK: supExpCK}
		dto.States = append(dto.States, st)
		return expCK
	case typedef.LinkRec:
		st := State{
			ExpCK:    expCK,
			K:        LinkKind,
			SupExpCK: supExpCK,
			Spec: ExpSpec{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typedef.PlusRec:
		st := State{
			ExpCK:    expCK,
			K:        PlusKind,
			SupExpCK: supExpCK,
			Spec: ExpSpec{Plus: &Sum{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpCK: statesFromExpRec(expCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typedef.WithRec:
		st := State{
			ExpCK:    expCK,
			K:        WithKind,
			SupExpCK: supExpCK,
			Spec: ExpSpec{With: &Sum{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpCK: statesFromExpRec(expCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typedef.UpRec:
		st := State{
			ExpCK:    expCK,
			K:        UpKind,
			SupExpCK: supExpCK,
			Spec: ExpSpec{Up: &Shift{
				ContExpCK: statesFromExpRec(expCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typedef.DownRec:
		st := State{
			ExpCK:    expCK,
			K:        DownKind,
			SupExpCK: supExpCK,
			Spec: ExpSpec{Down: &Shift{
				ContExpCK: statesFromExpRec(expCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expCK
	default:
		panic(typedef.ErrRecTypeUnexpected(r))
	}
}

func errExpKindUnexpected(k ExpKind) error {
	return fmt.Errorf("exp kind unexpected: %v", k)
}
