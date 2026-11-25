package adapter

import (
	"errors"
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	termexp "orglang/go-engine/internal/pool/termexp/core"
	compexp "orglang/go-engine/internal/pool/compexec/exp"
)

func DataFromCompSpec(s compexp.ExpSpec) CompSpec {
	switch spec := s.(type) {
	case compexp.AcquireSpec:
		return CompSpec{K: AcquireKind, Acquire: DataFromAcquireSpec(spec)}
	case compexp.AcceptSpec:
		return CompSpec{K: AcceptKind, Accept: DataFromAcceptSpec(spec)}
	case compexp.HireSpec:
		return CompSpec{K: HireKind, Hire: DataFromHireSpec(spec)}
	case compexp.ApplySpec:
		return CompSpec{K: ApplyKind, Apply: DataFromApplySpec(spec)}
	case compexp.ReleaseSpec:
		return CompSpec{K: ReleaseKind, Release: DataFromReleaseSpec(spec)}
	case compexp.DetachSpec:
		return CompSpec{K: DetachKind, Detach: DataFromDetachSpec(spec)}
	default:
		panic(compexp.ErrSpecTypeUnexpected(s))
	}
}

func DataToCompSpec(dto CompSpec) (compexp.ExpSpec, error) {
	switch dto.K {
	case AcquireKind:
		return DataToAcquireSpec(dto.Acquire)
	case AcceptKind:
		return DataToAcceptSpec(dto.Accept)
	case HireKind:
		return DataToHireSpec(dto.Hire)
	case ApplyKind:
		return DataToApplySpec(dto.Apply)
	case ReleaseKind:
		return DataToReleaseSpec(dto.Release)
	case DetachKind:
		return DataToDetachSpec(dto.Detach)
	default:
		panic(fmt.Errorf("exp kind unexpected: %v", dto.K))
	}
}

func DataFromExpRec(r compexp.ExpRec) ExpRec {
	switch rec := r.(type) {
	case compexp.AcquireRec:
		return ExpRec{K: AcquireKind, Acquire: DataFromAcquireRec(rec)}
	case compexp.AcceptRec:
		return ExpRec{K: AcceptKind, Accept: DataFromAcceptRec(rec)}
	case compexp.HireRec:
		return ExpRec{K: HireKind, Hire: DataFromHireRec(rec)}
	case compexp.ApplyRec:
		return ExpRec{K: ApplyKind, Apply: DataFromApplyRec(rec)}
	case compexp.ReleaseRec:
		return ExpRec{K: ReleaseKind, Release: DataFromReleaseRec(rec)}
	case compexp.DetachRec:
		return ExpRec{K: DetachKind, Detach: DataFromDetachRec(rec)}
	default:
		panic(compexp.ErrRecTypeUnexpected(r))
	}
}

func DataToExpRec(dto ExpRec) (compexp.ExpRec, error) {
	switch dto.K {
	case AcquireKind:
		return DataToAcquireRec(dto.Acquire)
	case AcceptKind:
		return DataToAcceptRec(dto.Accept)
	case HireKind:
		return DataToHireRec(dto.Hire)
	case ApplyKind:
		return DataToApplyRec(dto.Apply)
	case ReleaseKind:
		return DataToReleaseRec(dto.Release)
	case DetachKind:
		return DataToDetachRec(dto.Detach)
	default:
		panic(fmt.Errorf("exp kind unexpected: %v", dto.K))
	}
}

func dataFromDescRec(rec termexp.DescRec) (DescData, error) {
	if rec == nil {
		return DescData{}, errors.New("nil pool Description record")
	}
	dto := DescData{ExpCK: contkey.ConvertToInt(rec.ContentKey())}
	statesFromDescRec(0, rec, &dto)
	return dto, nil
}

func statesFromDescRec(parent int64, rec termexp.DescRec, dto *DescData) int64 {
	key := contkey.ConvertToInt(rec.ContentKey())
	var kind ExpKind
	data := StateSpec{}
	switch r := rec.(type) {
	case termexp.DescAcquireRec:
		kind, data.ConnChnlPH = AcquireKind, string(r.ConnChnlPH)
		data.ContExpCK = statesFromDescRec(key, r.ContExp, dto)
	case termexp.DescAcceptRec:
		kind, data.ConnChnlPH = AcceptKind, string(r.ConnChnlPH)
		data.ContExpCK = statesFromDescRec(key, r.ContExp, dto)
	case termexp.DescHireRec:
		kind, data.ConnChnlPH, data.ProcTermQN = HireKind, string(r.ConnChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
		data.ContExpCK = statesFromDescRec(key, r.ContExp, dto)
	case termexp.DescApplyRec:
		kind, data.ConnChnlPH, data.ProcTermQN = ApplyKind, string(r.ConnChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
		data.ContExpCK = statesFromDescRec(key, r.ContExp, dto)
	case termexp.DescReleaseRec:
		kind, data.ConnChnlPH = ReleaseKind, string(r.ConnChnlPH)
	case termexp.DescDetachRec:
		kind, data.ConnChnlPH = DetachKind, string(r.ConnChnlPH)
	case termexp.DescFireRec:
		kind, data.ConnChnlPH, data.ProcTermQN = termexp.FireKind, string(r.ConnChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
	case termexp.DescQuitRec:
		kind, data.ConnChnlPH, data.ProcTermQN = termexp.QuitKind, string(r.ConnChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
	case termexp.DescSpawnRec:
		kind, data.ProcTermQN = termexp.SpawnKind, uniqsym.ConvertToString(r.ProcTermQN)
		data.ProcCompQNs = uniqsym.ConvertToStrings(r.ProcCompQNs)
	default:
		panic(fmt.Sprintf("unexpected pool Description record %T", rec))
	}
	dto.States = append(dto.States, DescState{ExpCK: key, SupExpCK: parent, K: kind, Spec: data})
	return key
}

func dataToDescRec(dto DescData) (termexp.DescRec, error) {
	states := make(map[int64]DescState, len(dto.States))
	for _, state := range dto.States {
		states[state.ExpCK] = state
	}
	root, ok := states[dto.ExpCK]
	if !ok {
		return nil, fmt.Errorf("pool Description root %d is missing", dto.ExpCK)
	}
	return statesToDescRec(states, root)
}

func statesToDescRec(states map[int64]DescState, state DescState) (termexp.DescRec, error) {
	key, err := contkey.ConvertFromInt(state.ExpCK)
	if err != nil {
		return nil, err
	}
	spec := state.Spec
	cont := func() (termexp.DescRec, error) {
		child, ok := states[spec.ContExpCK]
		if !ok {
			return nil, fmt.Errorf("pool Description child %d is missing", spec.ContExpCK)
		}
		if child.SupExpCK != state.ExpCK {
			return nil, fmt.Errorf("pool Description child %d has unexpected parent", child.ExpCK)
		}
		return statesToDescRec(states, child)
	}
	ph := func() (symbol.ADT, error) { return symbol.ConvertFromString(spec.ConnChnlPH) }
	switch state.K {
	case AcquireKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		c, e := cont()
		return termexp.DescAcquireRec{ExpCK: key, ConnChnlPH: p, ContExp: c}, e
	case AcceptKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		c, e := cont()
		return termexp.DescAcceptRec{ExpCK: key, ConnChnlPH: p, ContExp: c}, e
	case HireKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		if e != nil {
			return nil, e
		}
		c, e := cont()
		return termexp.DescHireRec{ExpCK: key, ConnChnlPH: p, ProcTermQN: q, ContExp: c}, e
	case ApplyKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		if e != nil {
			return nil, e
		}
		c, e := cont()
		return termexp.DescApplyRec{ExpCK: key, ConnChnlPH: p, ProcTermQN: q, ContExp: c}, e
	case ReleaseKind:
		p, e := ph()
		return termexp.DescReleaseRec{ExpCK: key, ConnChnlPH: p}, e
	case DetachKind:
		p, e := ph()
		return termexp.DescDetachRec{ExpCK: key, ConnChnlPH: p}, e
	case termexp.FireKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		return termexp.DescFireRec{ExpCK: key, ConnChnlPH: p, ProcTermQN: q}, e
	case termexp.QuitKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		return termexp.DescQuitRec{ExpCK: key, ConnChnlPH: p, ProcTermQN: q}, e
	case termexp.SpawnKind:
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		if e != nil {
			return nil, e
		}
		qs, e := uniqsym.ConvertFromStrings(spec.ProcCompQNs)
		return termexp.DescSpawnRec{ExpCK: key, ProcTermQN: q, ProcCompQNs: qs}, e
	default:
		return nil, fmt.Errorf("unexpected pool Description kind %v", state.K)
	}
}
