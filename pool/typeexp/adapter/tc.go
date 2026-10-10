package adapter

import (
	"fmt"

	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"

	"github.com/orglang/go-sdk/pool/typeexp"

	pooltypeexp "orglang/go-engine/pool/typeexp/core"
)

func MsgFromExpSpec(s pooltypeexp.ExpSpec) typeexp.ExpSpec {
	switch spec := s.(type) {
	case pooltypeexp.OneSpec:
		return typeexp.ExpSpec{K: typeexp.One}
	case pooltypeexp.LinkSpec:
		return typeexp.ExpSpec{
			K:    typeexp.Link,
			Link: &typeexp.LinkSpec{TypeQN: uniqsym.ConvertToString(spec.TypeQN)}}
	case pooltypeexp.WithSpec:
		return typeexp.ExpSpec{
			K: typeexp.With,
			With: &typeexp.LaborSpec{
				ProcQNs: uniqsym.ConvertToStrings(spec.ProcQNs),
				ContExp: MsgFromExpSpec(spec.ContExp)},
		}
	case pooltypeexp.PlusSpec:
		return typeexp.ExpSpec{
			K: typeexp.Plus,
			Plus: &typeexp.LaborSpec{
				ProcQNs: uniqsym.ConvertToStrings(spec.ProcQNs),
				ContExp: MsgFromExpSpec(spec.ContExp)},
		}
	default:
		panic(pooltypeexp.ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpec(dto typeexp.ExpSpec) (pooltypeexp.ExpSpec, error) {
	switch dto.K {
	case typeexp.One:
		return pooltypeexp.OneSpec{}, nil
	case typeexp.Link:
		xactQN, err := uniqsym.ConvertFromString(dto.Link.TypeQN)
		if err != nil {
			return nil, err
		}
		return pooltypeexp.LinkSpec{TypeQN: xactQN}, nil
	case typeexp.Plus:
		procQNs, err := uniqsym.ConvertFromStrings(dto.Plus.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := MsgToExpSpec(dto.Plus.ContExp)
		if err != nil {
			return nil, err
		}
		return pooltypeexp.PlusSpec{ProcQNs: procQNs, ContExp: contExp}, nil
	case typeexp.With:
		procQNs, err := uniqsym.ConvertFromStrings(dto.With.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := MsgToExpSpec(dto.With.ContExp)
		if err != nil {
			return nil, err
		}
		return pooltypeexp.WithSpec{ProcQNs: procQNs, ContExp: contExp}, nil
	case typeexp.Up:
		contExp, err := MsgToExpSpec(dto.Up.ContExp)
		if err != nil {
			return nil, err
		}
		return pooltypeexp.UpSpec{ContExp: contExp}, nil
	case typeexp.Down:
		contExp, err := MsgToExpSpec(dto.Down.ContExp)
		if err != nil {
			return nil, err
		}
		return pooltypeexp.DownSpec{ContExp: contExp}, nil
	default:
		panic(typeexp.ErrKindUnexpected(dto.K))
	}
}

func msgFromExpRef(ref pooltypeexp.ExpRef) typeexp.ExpRef {
	expVK := valkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case pooltypeexp.OneRef, pooltypeexp.OneRec:
		return typeexp.ExpRef{K: typeexp.One, ExpVK: expVK}
	case pooltypeexp.LinkRef, pooltypeexp.LinkRec:
		return typeexp.ExpRef{K: typeexp.Link, ExpVK: expVK}
	case pooltypeexp.PlusRef, pooltypeexp.PlusRec:
		return typeexp.ExpRef{K: typeexp.Plus, ExpVK: expVK}
	case pooltypeexp.WithRef, pooltypeexp.WithRec:
		return typeexp.ExpRef{K: typeexp.With, ExpVK: expVK}
	default:
		panic(pooltypeexp.ErrRefTypeUnexpected(ref))
	}
}

func msgToExpRef(dto typeexp.ExpRef) (pooltypeexp.ExpRef, error) {
	expVK, err := valkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case typeexp.One:
		return pooltypeexp.OneRef{ExpVK: expVK}, nil
	case typeexp.Link:
		return pooltypeexp.LinkRef{ExpVK: expVK}, nil
	case typeexp.Plus:
		return pooltypeexp.PlusRef{ExpVK: expVK}, nil
	case typeexp.With:
		return pooltypeexp.WithRef{ExpVK: expVK}, nil
	default:
		panic(typeexp.ErrKindUnexpected(dto.K))
	}
}

func dataFromExpRef(ref pooltypeexp.ExpRef) pooltypeexp.ExpRefDS {
	if ref == nil {
		panic("can't be nil")
	}
	expVK := valkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case pooltypeexp.OneRef, pooltypeexp.OneRec:
		return pooltypeexp.ExpRefDS{K: pooltypeexp.OneKind, ExpVK: expVK}
	case pooltypeexp.LinkRef, pooltypeexp.LinkRec:
		return pooltypeexp.ExpRefDS{K: pooltypeexp.LinkKind, ExpVK: expVK}
	case pooltypeexp.PlusRef, pooltypeexp.PlusRec:
		return pooltypeexp.ExpRefDS{K: pooltypeexp.PlusKind, ExpVK: expVK}
	case pooltypeexp.WithRef, pooltypeexp.WithRec:
		return pooltypeexp.ExpRefDS{K: pooltypeexp.WithKind, ExpVK: expVK}
	default:
		panic(pooltypeexp.ErrRefTypeUnexpected(ref))
	}
}

func dataToExpRef(dto pooltypeexp.ExpRefDS) (pooltypeexp.ExpRef, error) {
	expVK, err := valkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case pooltypeexp.OneKind:
		return pooltypeexp.OneRef{ExpVK: expVK}, nil
	case pooltypeexp.LinkKind:
		return pooltypeexp.LinkRef{ExpVK: expVK}, nil
	case pooltypeexp.PlusKind:
		return pooltypeexp.PlusRef{ExpVK: expVK}, nil
	case pooltypeexp.WithKind:
		return pooltypeexp.WithRef{ExpVK: expVK}, nil
	default:
		panic(errExpKindUnexpected(dto.K))
	}
}

func dataToExpRec(dto pooltypeexp.ExpRecDS) (pooltypeexp.ExpRec, error) {
	states := make(map[int64]pooltypeexp.StateDS, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ExpVK] = dto
	}
	return statesToExpRec(states, states[dto.ExpVK])
}

func dataFromExpRec(rec pooltypeexp.ExpRec) pooltypeexp.ExpRecDS {
	if rec == nil {
		panic("can't be nil")
	}
	dto := &pooltypeexp.ExpRecDS{
		ExpVK:  valkey.ConvertToInt(rec.Key()),
		States: nil,
	}
	statesFromExpRec(0, rec, dto)
	return *dto
}

func statesToExpRec(states map[int64]pooltypeexp.StateDS, st pooltypeexp.StateDS) (pooltypeexp.ExpRec, error) {
	expVK, err := valkey.ConvertFromInt(st.ExpVK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case pooltypeexp.OneKind:
		return pooltypeexp.OneRec{ExpVK: expVK}, nil
	case pooltypeexp.LinkKind:
		xactQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return pooltypeexp.LinkRec{ExpVK: expVK, TypeQN: xactQN}, nil
	case pooltypeexp.PlusKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.Plus.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := statesToExpRec(states, states[st.Spec.Plus.ContExpVK])
		if err != nil {
			return nil, err
		}
		return pooltypeexp.PlusRec{ExpVK: expVK, ProcQNs: procQNs, ContExp: contExp}, nil
	case pooltypeexp.WithKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.With.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := statesToExpRec(states, states[st.Spec.With.ContExpVK])
		if err != nil {
			return nil, err
		}
		return pooltypeexp.WithRec{ExpVK: expVK, ProcQNs: procQNs, ContExp: contExp}, nil
	case pooltypeexp.UpKind:
		contExp, err := statesToExpRec(states, states[st.Spec.Up.ContExpVK])
		if err != nil {
			return nil, err
		}
		return pooltypeexp.UpRec{ExpVK: expVK, ContExp: contExp}, nil
	case pooltypeexp.DownKind:
		contExp, err := statesToExpRec(states, states[st.Spec.Down.ContExpVK])
		if err != nil {
			return nil, err
		}
		return pooltypeexp.DownRec{ExpVK: expVK, ContExp: contExp}, nil
	default:
		panic(errExpKindUnexpected(st.K))
	}
}

func statesFromExpRec(supExpVK int64, r pooltypeexp.ExpRec, dto *pooltypeexp.ExpRecDS) int64 {
	expVK := valkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case pooltypeexp.OneRec:
		st := pooltypeexp.StateDS{ExpVK: expVK, K: pooltypeexp.OneKind, SupExpVK: supExpVK}
		dto.States = append(dto.States, st)
		return expVK
	case pooltypeexp.LinkRec:
		st := pooltypeexp.StateDS{
			ExpVK:    expVK,
			K:        pooltypeexp.LinkKind,
			SupExpVK: supExpVK,
			Spec: pooltypeexp.ExpSpecDS{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case pooltypeexp.PlusRec:
		st := pooltypeexp.StateDS{
			ExpVK:    expVK,
			K:        pooltypeexp.PlusKind,
			SupExpVK: supExpVK,
			Spec: pooltypeexp.ExpSpecDS{Plus: &pooltypeexp.SumDS{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	case pooltypeexp.WithRec:
		st := pooltypeexp.StateDS{
			ExpVK:    expVK,
			K:        pooltypeexp.WithKind,
			SupExpVK: supExpVK,
			Spec: pooltypeexp.ExpSpecDS{With: &pooltypeexp.SumDS{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	case pooltypeexp.UpRec:
		st := pooltypeexp.StateDS{
			ExpVK:    expVK,
			K:        pooltypeexp.UpKind,
			SupExpVK: supExpVK,
			Spec: pooltypeexp.ExpSpecDS{Up: &pooltypeexp.ShiftDS{
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	case pooltypeexp.DownRec:
		st := pooltypeexp.StateDS{
			ExpVK:    expVK,
			K:        pooltypeexp.DownKind,
			SupExpVK: supExpVK,
			Spec: pooltypeexp.ExpSpecDS{Up: &pooltypeexp.ShiftDS{
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	default:
		panic(pooltypeexp.ErrRecTypeUnexpected(r))
	}
}

func errExpKindUnexpected(k pooltypeexp.ExpKind) error {
	return fmt.Errorf("exp kind unexpected: %v", k)
}
