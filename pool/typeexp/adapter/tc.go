package adapter

import (
	"fmt"

	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"

	sdktypeexp "github.com/orglang/go-sdk/pool/typeexp"

	typeexp "orglang/go-engine/pool/typeexp/core"
)

func MsgFromExpSpec(s typeexp.ExpSpec) sdktypeexp.ExpSpec {
	switch spec := s.(type) {
	case typeexp.OneSpec:
		return sdktypeexp.ExpSpec{K: sdktypeexp.One}
	case typeexp.LinkSpec:
		return sdktypeexp.ExpSpec{
			K:    sdktypeexp.Link,
			Link: &sdktypeexp.LinkSpec{TypeQN: uniqsym.ConvertToString(spec.TypeQN)}}
	case typeexp.WithSpec:
		return sdktypeexp.ExpSpec{
			K: sdktypeexp.With,
			With: &sdktypeexp.LaborSpec{
				ProcQNs: uniqsym.ConvertToStrings(spec.ProcQNs),
				ContExp: MsgFromExpSpec(spec.ContExp)},
		}
	case typeexp.PlusSpec:
		return sdktypeexp.ExpSpec{
			K: sdktypeexp.Plus,
			Plus: &sdktypeexp.LaborSpec{
				ProcQNs: uniqsym.ConvertToStrings(spec.ProcQNs),
				ContExp: MsgFromExpSpec(spec.ContExp)},
		}
	default:
		panic(typeexp.ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpec(dto sdktypeexp.ExpSpec) (typeexp.ExpSpec, error) {
	switch dto.K {
	case sdktypeexp.One:
		return typeexp.OneSpec{}, nil
	case sdktypeexp.Link:
		xactQN, err := uniqsym.ConvertFromString(dto.Link.TypeQN)
		if err != nil {
			return nil, err
		}
		return typeexp.LinkSpec{TypeQN: xactQN}, nil
	case sdktypeexp.Plus:
		procQNs, err := uniqsym.ConvertFromStrings(dto.Plus.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := MsgToExpSpec(dto.Plus.ContExp)
		if err != nil {
			return nil, err
		}
		return typeexp.PlusSpec{ProcQNs: procQNs, ContExp: contExp}, nil
	case sdktypeexp.With:
		procQNs, err := uniqsym.ConvertFromStrings(dto.With.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := MsgToExpSpec(dto.With.ContExp)
		if err != nil {
			return nil, err
		}
		return typeexp.WithSpec{ProcQNs: procQNs, ContExp: contExp}, nil
	case sdktypeexp.Up:
		contExp, err := MsgToExpSpec(dto.Up.ContExp)
		if err != nil {
			return nil, err
		}
		return typeexp.UpSpec{ContExp: contExp}, nil
	case sdktypeexp.Down:
		contExp, err := MsgToExpSpec(dto.Down.ContExp)
		if err != nil {
			return nil, err
		}
		return typeexp.DownSpec{ContExp: contExp}, nil
	default:
		panic(sdktypeexp.ErrKindUnexpected(dto.K))
	}
}

func msgFromExpRef(ref typeexp.ExpRef) sdktypeexp.ExpRef {
	expVK := valkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typeexp.OneRef, typeexp.OneRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.One, ExpVK: expVK}
	case typeexp.LinkRef, typeexp.LinkRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.Link, ExpVK: expVK}
	case typeexp.PlusRef, typeexp.PlusRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.Plus, ExpVK: expVK}
	case typeexp.WithRef, typeexp.WithRec:
		return sdktypeexp.ExpRef{K: sdktypeexp.With, ExpVK: expVK}
	default:
		panic(typeexp.ErrRefTypeUnexpected(ref))
	}
}

func msgToExpRef(dto sdktypeexp.ExpRef) (typeexp.ExpRef, error) {
	expVK, err := valkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case sdktypeexp.One:
		return typeexp.OneRef{ExpVK: expVK}, nil
	case sdktypeexp.Link:
		return typeexp.LinkRef{ExpVK: expVK}, nil
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
		panic("can't be nil")
	}
	expVK := valkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case typeexp.OneRef, typeexp.OneRec:
		return typeexp.ExpRefDS{K: typeexp.OneKind, ExpVK: expVK}
	case typeexp.LinkRef, typeexp.LinkRec:
		return typeexp.ExpRefDS{K: typeexp.LinkKind, ExpVK: expVK}
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
	case typeexp.PlusKind:
		return typeexp.PlusRef{ExpVK: expVK}, nil
	case typeexp.WithKind:
		return typeexp.WithRef{ExpVK: expVK}, nil
	default:
		panic(errExpKindUnexpected(dto.K))
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
		xactQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return typeexp.LinkRec{ExpVK: expVK, TypeQN: xactQN}, nil
	case typeexp.PlusKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.Plus.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := statesToExpRec(states, states[st.Spec.Plus.ContExpVK])
		if err != nil {
			return nil, err
		}
		return typeexp.PlusRec{ExpVK: expVK, ProcQNs: procQNs, ContExp: contExp}, nil
	case typeexp.WithKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.With.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := statesToExpRec(states, states[st.Spec.With.ContExpVK])
		if err != nil {
			return nil, err
		}
		return typeexp.WithRec{ExpVK: expVK, ProcQNs: procQNs, ContExp: contExp}, nil
	case typeexp.UpKind:
		contExp, err := statesToExpRec(states, states[st.Spec.Up.ContExpVK])
		if err != nil {
			return nil, err
		}
		return typeexp.UpRec{ExpVK: expVK, ContExp: contExp}, nil
	case typeexp.DownKind:
		contExp, err := statesToExpRec(states, states[st.Spec.Down.ContExpVK])
		if err != nil {
			return nil, err
		}
		return typeexp.DownRec{ExpVK: expVK, ContExp: contExp}, nil
	default:
		panic(errExpKindUnexpected(st.K))
	}
}

func statesFromExpRec(supExpVK int64, r typeexp.ExpRec, dto *typeexp.ExpRecDS) int64 {
	expVK := valkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case typeexp.OneRec:
		st := typeexp.StateDS{ExpVK: expVK, K: typeexp.OneKind, SupExpVK: supExpVK}
		dto.States = append(dto.States, st)
		return expVK
	case typeexp.LinkRec:
		st := typeexp.StateDS{
			ExpVK:    expVK,
			K:        typeexp.LinkKind,
			SupExpVK: supExpVK,
			Spec: typeexp.ExpSpecDS{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case typeexp.PlusRec:
		st := typeexp.StateDS{
			ExpVK:    expVK,
			K:        typeexp.PlusKind,
			SupExpVK: supExpVK,
			Spec: typeexp.ExpSpecDS{Plus: &typeexp.SumDS{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	case typeexp.WithRec:
		st := typeexp.StateDS{
			ExpVK:    expVK,
			K:        typeexp.WithKind,
			SupExpVK: supExpVK,
			Spec: typeexp.ExpSpecDS{With: &typeexp.SumDS{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	case typeexp.UpRec:
		st := typeexp.StateDS{
			ExpVK:    expVK,
			K:        typeexp.UpKind,
			SupExpVK: supExpVK,
			Spec: typeexp.ExpSpecDS{Up: &typeexp.ShiftDS{
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	case typeexp.DownRec:
		st := typeexp.StateDS{
			ExpVK:    expVK,
			K:        typeexp.DownKind,
			SupExpVK: supExpVK,
			Spec: typeexp.ExpSpecDS{Up: &typeexp.ShiftDS{
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	default:
		panic(typeexp.ErrRecTypeUnexpected(r))
	}
}

func errExpKindUnexpected(k typeexp.ExpKind) error {
	return fmt.Errorf("exp kind unexpected: %v", k)
}
