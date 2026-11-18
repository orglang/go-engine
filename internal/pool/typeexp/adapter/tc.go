package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"

	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

func dataFromExpRef(ref typeexp.ExpRef) typeexp.ExpRefDS {
	if ref == nil {
		panic("can't be nil")
	}
	expCK := contkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typeexp.OneRef, typeexp.OneRec:
		return typeexp.ExpRefDS{K: typeexp.OneKind, ExpCK: expCK}
	case typeexp.LinkRef, typeexp.LinkRec:
		return typeexp.ExpRefDS{K: typeexp.LinkKind, ExpCK: expCK}
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
	case typeexp.PlusKind:
		return typeexp.PlusRef{ExpCK: expCK}, nil
	case typeexp.WithKind:
		return typeexp.WithRef{ExpCK: expCK}, nil
	default:
		panic(errExpKindUnexpected(dto.K))
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
		xactQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return typeexp.LinkRec{ExpCK: expCK, TypeQN: xactQN}, nil
	case typeexp.PlusKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.Plus.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := statesToExpRec(states, states[st.Spec.Plus.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typeexp.PlusRec{ExpCK: expCK, ProcQNs: procQNs, ContExp: contExp}, nil
	case typeexp.WithKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.With.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := statesToExpRec(states, states[st.Spec.With.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typeexp.WithRec{ExpCK: expCK, ProcQNs: procQNs, ContExp: contExp}, nil
	case typeexp.UpKind:
		contExp, err := statesToExpRec(states, states[st.Spec.Up.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typeexp.UpRec{ExpCK: expCK, ContExp: contExp}, nil
	case typeexp.DownKind:
		contExp, err := statesToExpRec(states, states[st.Spec.Down.ContExpCK])
		if err != nil {
			return nil, err
		}
		return typeexp.DownRec{ExpCK: expCK, ContExp: contExp}, nil
	default:
		panic(errExpKindUnexpected(st.K))
	}
}

func statesFromExpRec(supExpCK int64, r typeexp.ExpRec, dto *typeexp.ExpRecDS) int64 {
	expCK := contkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case typeexp.OneRec:
		st := typeexp.StateDS{ExpCK: expCK, K: typeexp.OneKind, SupExpCK: supExpCK}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.LinkRec:
		st := typeexp.StateDS{
			ExpCK:    expCK,
			K:        typeexp.LinkKind,
			SupExpCK: supExpCK,
			Spec: typeexp.ExpSpecDS{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.PlusRec:
		st := typeexp.StateDS{
			ExpCK:    expCK,
			K:        typeexp.PlusKind,
			SupExpCK: supExpCK,
			Spec: typeexp.ExpSpecDS{Plus: &typeexp.SumDS{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpCK: statesFromExpRec(expCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.WithRec:
		st := typeexp.StateDS{
			ExpCK:    expCK,
			K:        typeexp.WithKind,
			SupExpCK: supExpCK,
			Spec: typeexp.ExpSpecDS{With: &typeexp.SumDS{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpCK: statesFromExpRec(expCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.UpRec:
		st := typeexp.StateDS{
			ExpCK:    expCK,
			K:        typeexp.UpKind,
			SupExpCK: supExpCK,
			Spec: typeexp.ExpSpecDS{Up: &typeexp.ShiftDS{
				ContExpCK: statesFromExpRec(expCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expCK
	case typeexp.DownRec:
		st := typeexp.StateDS{
			ExpCK:    expCK,
			K:        typeexp.DownKind,
			SupExpCK: supExpCK,
			Spec: typeexp.ExpSpecDS{Down: &typeexp.ShiftDS{
				ContExpCK: statesFromExpRec(expCK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expCK
	default:
		panic(typeexp.ErrRecTypeUnexpected(r))
	}
}

func errExpKindUnexpected(k typeexp.ExpKind) error {
	return fmt.Errorf("exp kind unexpected: %v", k)
}
