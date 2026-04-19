package typeexp

import (
	"fmt"

	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"

	"github.com/orglang/go-sdk/pool/typeexp"
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
	case WithSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		vks := make([]valkey.ADT, 0, len(spec.ProcQNs)+1)
		for _, qn := range spec.ProcQNs {
			vk, err := qn.Key()
			if err != nil {
				return nil, err
			}
			vks = append(vks, vk)
		}
		expVK, err := valkey.Compose(append(vks, contExp.Key())...)
		if err != nil {
			return nil, err
		}
		return WithRec{ExpVK: expVK, ProcQNs: spec.ProcQNs, ContExp: contExp}, nil
	case PlusSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		vks := make([]valkey.ADT, 0, len(spec.ProcQNs)+1)
		for _, qn := range spec.ProcQNs {
			vk, err := qn.Key()
			if err != nil {
				return nil, err
			}
			vks = append(vks, vk)
		}
		expVK, err := valkey.Compose(append(vks, contExp.Key())...)
		return PlusRec{ExpVK: expVK, ProcQNs: spec.ProcQNs, ContExp: contExp}, nil
	case UpSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		expVK, err := valkey.Compose([]valkey.ADT{valkey.Two, contExp.Key()}...)
		return UpRec{ExpVK: expVK, ContExp: contExp}, nil
	case DownSpec:
		contExp, err := ConvertSpecToRec(spec.ContExp)
		if err != nil {
			return nil, err
		}
		expVK, err := valkey.Compose([]valkey.ADT{valkey.Three, contExp.Key()}...)
		return UpRec{ExpVK: expVK, ContExp: contExp}, nil
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
	case WithRec:
		return WithSpec{ProcQNs: rec.ProcQNs, ContExp: ConvertRecToSpec(rec.ContExp)}
	case PlusRec:
		return PlusSpec{ProcQNs: rec.ProcQNs, ContExp: ConvertRecToSpec(rec.ContExp)}
	case UpRec:
		return UpSpec{ContExp: ConvertRecToSpec(rec.ContExp)}
	case DownRec:
		return DownSpec{ContExp: ConvertRecToSpec(rec.ContExp)}
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
			Link: &typeexp.LinkSpec{XactQN: uniqsym.ConvertToString(spec.TypeQN)}}
	case WithSpec:
		return typeexp.ExpSpec{
			K: typeexp.With,
			With: &typeexp.LaborSpec{
				ProcQNs: uniqsym.ConvertToStrings(spec.ProcQNs),
				ContExp: MsgFromExpSpec(spec.ContExp)},
		}
	case PlusSpec:
		return typeexp.ExpSpec{
			K: typeexp.Plus,
			Plus: &typeexp.LaborSpec{
				ProcQNs: uniqsym.ConvertToStrings(spec.ProcQNs),
				ContExp: MsgFromExpSpec(spec.ContExp)},
		}
	default:
		panic(ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpec(dto typeexp.ExpSpec) (ExpSpec, error) {
	switch dto.K {
	case typeexp.One:
		return OneSpec{}, nil
	case typeexp.Link:
		xactQN, err := uniqsym.ConvertFromString(dto.Link.XactQN)
		if err != nil {
			return nil, err
		}
		return LinkSpec{TypeQN: xactQN}, nil
	case typeexp.Plus:
		procQNs, err := uniqsym.ConvertFromStrings(dto.Plus.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := MsgToExpSpec(dto.Plus.ContExp)
		if err != nil {
			return nil, err
		}
		return PlusSpec{ProcQNs: procQNs, ContExp: contExp}, nil
	case typeexp.With:
		procQNs, err := uniqsym.ConvertFromStrings(dto.With.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := MsgToExpSpec(dto.With.ContExp)
		if err != nil {
			return nil, err
		}
		return WithSpec{ProcQNs: procQNs, ContExp: contExp}, nil
	case typeexp.Up:
		contExp, err := MsgToExpSpec(dto.Up.ContExp)
		if err != nil {
			return nil, err
		}
		return UpSpec{ContExp: contExp}, nil
	case typeexp.Down:
		contExp, err := MsgToExpSpec(dto.Down.ContExp)
		if err != nil {
			return nil, err
		}
		return DownSpec{ContExp: contExp}, nil
	default:
		panic(typeexp.ErrKindUnexpected(dto.K))
	}
}

func msgFromExpRef(ref ExpRef) typeexp.ExpRef {
	expVK := valkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case OneRef, OneRec:
		return typeexp.ExpRef{K: typeexp.One, ExpVK: expVK}
	case LinkRef, LinkRec:
		return typeexp.ExpRef{K: typeexp.Link, ExpVK: expVK}
	case PlusRef, PlusRec:
		return typeexp.ExpRef{K: typeexp.Plus, ExpVK: expVK}
	case WithRef, WithRec:
		return typeexp.ExpRef{K: typeexp.With, ExpVK: expVK}
	default:
		panic(ErrRefTypeUnexpected(ref))
	}
}

func msgToExpRef(dto typeexp.ExpRef) (ExpRef, error) {
	expVK, err := valkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case typeexp.One:
		return OneRef{expVK}, nil
	case typeexp.Link:
		return LinkRef{expVK}, nil
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
		panic("can't be nil")
	}
	expVK := valkey.ConvertToInt(ref.Key())
	switch ref.(type) {
	case OneRef, OneRec:
		return expRefDS{K: oneKind, ExpVK: expVK}
	case LinkRef, LinkRec:
		return expRefDS{K: linkKind, ExpVK: expVK}
	case PlusRef, PlusRec:
		return expRefDS{K: plusKind, ExpVK: expVK}
	case WithRef, WithRec:
		return expRefDS{K: withKind, ExpVK: expVK}
	default:
		panic(ErrRefTypeUnexpected(ref))
	}
}

func dataToExpRef(dto expRefDS) (ExpRef, error) {
	expVK, err := valkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case oneKind:
		return OneRef{expVK}, nil
	case linkKind:
		return LinkRef{expVK}, nil
	case plusKind:
		return PlusRef{expVK}, nil
	case withKind:
		return WithRef{expVK}, nil
	default:
		panic(errExpKindUnexpected(dto.K))
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
		ExpVK:  valkey.ConvertToInt(rec.Key()),
		States: nil,
	}
	statesFromExpRec(0, rec, dto)
	return *dto
}

func statesToExpRec(states map[int64]stateDS, st stateDS) (ExpRec, error) {
	expVK, err := valkey.ConvertFromInt(st.ExpVK)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case oneKind:
		return OneRec{ExpVK: expVK}, nil
	case linkKind:
		xactQN, err := uniqsym.ConvertFromString(st.Spec.Link)
		if err != nil {
			return nil, err
		}
		return LinkRec{ExpVK: expVK, TypeQN: xactQN}, nil
	case plusKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.Plus.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := statesToExpRec(states, states[st.Spec.Plus.ContExpVK])
		if err != nil {
			return nil, err
		}
		return PlusRec{ExpVK: expVK, ProcQNs: procQNs, ContExp: contExp}, nil
	case withKind:
		procQNs, err := uniqsym.ConvertFromStrings(st.Spec.With.ProcQNs)
		if err != nil {
			return nil, err
		}
		contExp, err := statesToExpRec(states, states[st.Spec.With.ContExpVK])
		if err != nil {
			return nil, err
		}
		return WithRec{ExpVK: expVK, ProcQNs: procQNs, ContExp: contExp}, nil
	case upKind:
		contExp, err := statesToExpRec(states, states[st.Spec.Up.ContExpVK])
		if err != nil {
			return nil, err
		}
		return UpRec{ExpVK: expVK, ContExp: contExp}, nil
	case downKind:
		contExp, err := statesToExpRec(states, states[st.Spec.Down.ContExpVK])
		if err != nil {
			return nil, err
		}
		return DownRec{ExpVK: expVK, ContExp: contExp}, nil
	default:
		panic(errExpKindUnexpected(st.K))
	}
}

func statesFromExpRec(supExpVK int64, r ExpRec, dto *expRecDS) int64 {
	expVK := valkey.ConvertToInt(r.Key())
	switch rec := r.(type) {
	case OneRec:
		st := stateDS{ExpVK: expVK, K: oneKind, SupExpVK: supExpVK}
		dto.States = append(dto.States, st)
		return expVK
	case LinkRec:
		st := stateDS{
			ExpVK:    expVK,
			K:        linkKind,
			SupExpVK: supExpVK,
			Spec: expSpecDS{
				Link: uniqsym.ConvertToString(rec.TypeQN),
			},
		}
		dto.States = append(dto.States, st)
		return expVK
	case PlusRec:
		st := stateDS{
			ExpVK:    expVK,
			K:        plusKind,
			SupExpVK: supExpVK,
			Spec: expSpecDS{Plus: &sumDS{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	case WithRec:
		st := stateDS{
			ExpVK:    expVK,
			K:        withKind,
			SupExpVK: supExpVK,
			Spec: expSpecDS{With: &sumDS{
				ProcQNs:   uniqsym.ConvertToStrings(rec.ProcQNs),
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	case UpRec:
		st := stateDS{
			ExpVK:    expVK,
			K:        upKind,
			SupExpVK: supExpVK,
			Spec: expSpecDS{Up: &shiftDS{
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	case DownRec:
		st := stateDS{
			ExpVK:    expVK,
			K:        downKind,
			SupExpVK: supExpVK,
			Spec: expSpecDS{Up: &shiftDS{
				ContExpVK: statesFromExpRec(expVK, rec.ContExp, dto),
			}},
		}
		dto.States = append(dto.States, st)
		return expVK
	default:
		panic(ErrRecTypeUnexpected(r))
	}
}

func errExpKindUnexpected(k expKind) error {
	return fmt.Errorf("exp kind unexpected: %v", k)
}
