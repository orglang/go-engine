package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"

	"orglang/go-engine/internal/adt/uniqsym"

	termexp "orglang/go-engine/internal/pool/termexp/core"
)

func DataFromExpSpec(s termexp.CompSpec) ExpSpec {
	switch spec := s.(type) {
	case termexp.AcquireSpec:
		return ExpSpec{K: AcquireKind, Acquire: DataFromAcquireSpec(spec)}
	case termexp.AcceptSpec:
		return ExpSpec{K: AcceptKind, Accept: DataFromAcceptSpec(spec)}
	case termexp.HireSpec:
		return ExpSpec{K: HireKind, Hire: DataFromHireSpec(spec)}
	case termexp.ApplySpec:
		return ExpSpec{K: ApplyKind, Apply: DataFromApplySpec(spec)}
	case termexp.ReleaseSpec:
		return ExpSpec{K: ReleaseKind, Release: DataFromReleaseSpec(spec)}
	case termexp.DetachSpec:
		return ExpSpec{K: DetachKind, Detach: DataFromDetachSpec(spec)}
	default:
		panic(termexp.ErrSpecTypeUnexpected(s))
	}
}

func DataToExpSpec(dto ExpSpec) (termexp.CompSpec, error) {
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

func DataFromExpRec(r termexp.CompRec) ExpRec {
	switch rec := r.(type) {
	case termexp.AcquireRec:
		return ExpRec{K: AcquireKind, Acquire: DataFromAcquireRec(rec)}
	case termexp.AcceptRec:
		return ExpRec{K: AcceptKind, Accept: DataFromAcceptRec(rec)}
	case termexp.HireRec:
		return ExpRec{K: HireKind, Hire: DataFromHireRec(rec)}
	case termexp.ApplyRec:
		return ExpRec{K: ApplyKind, Apply: DataFromApplyRec(rec)}
	case termexp.ReleaseRec:
		return ExpRec{K: ReleaseKind, Release: DataFromReleaseRec(rec)}
	case termexp.DetachRec:
		return ExpRec{K: DetachKind, Detach: DataFromDetachRec(rec)}
	default:
		panic(termexp.ErrRecTypeUnexpected(r))
	}
}

func DataToExpRec(dto ExpRec) (termexp.CompRec, error) {
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
		return DescData{}, fmt.Errorf("nil pool Description record")
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
		kind, data.ConnChnlPH, data.ProcTermQN = ExpKind(termexp.FireKind), string(r.ConnChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
	case termexp.DescQuitRec:
		kind, data.ConnChnlPH, data.ProcTermQN = ExpKind(termexp.QuitKind), string(r.ConnChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
	case termexp.DescSpawnRec:
		kind, data.ProcTermQN = ExpKind(termexp.SpawnKind), uniqsym.ConvertToString(r.ProcTermQN)
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
	case ExpKind(termexp.FireKind):
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		return termexp.DescFireRec{ExpCK: key, ConnChnlPH: p, ProcTermQN: q}, e
	case ExpKind(termexp.QuitKind):
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		return termexp.DescQuitRec{ExpCK: key, ConnChnlPH: p, ProcTermQN: q}, e
	case ExpKind(termexp.SpawnKind):
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
