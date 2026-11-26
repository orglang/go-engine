package adapter

import (
	"errors"
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

const (
	AcquireKind int16 = iota + 1
	AcceptKind
	HireKind
	ApplyKind
	ReleaseKind
	DetachKind
	FireKind
	QuitKind
	SpawnKind
)

type ExpData struct {
	ExpCK  int64      `db:"exp_ck"`
	States []ExpState `db:"states"`
}

type ExpState struct {
	ExpCK    int64        `db:"exp_ck"`
	SupExpCK int64        `db:"sup_exp_ck"`
	K        int16        `db:"kind"`
	Spec     ExpStateSpec `db:"spec"`
}

type ExpStateSpec struct {
	ConnChnlPH  string   `json:"conn_chnl_ph,omitempty"`
	ProcTermQN  string   `json:"proc_term_qn,omitempty"`
	ProcCompQNs []string `json:"proc_comp_qns,omitempty"`
	ContExpCK   int64    `json:"cont_exp_ck,omitempty"`
}

func dataFromExpRec(rec termdef.ExpRec) (ExpData, error) {
	if rec == nil {
		return ExpData{}, errors.New("nil pool Definition record")
	}
	dto := ExpData{ExpCK: contkey.ConvertToInt(rec.ContentKey())}
	statesFromExpRec(0, rec, &dto)
	return dto, nil
}

func statesFromExpRec(parent int64, rec termdef.ExpRec, dto *ExpData) int64 {
	key := contkey.ConvertToInt(rec.ContentKey())
	var kind int16
	data := ExpStateSpec{}
	switch r := rec.(type) {
	case termdef.ExpAcquireRec:
		kind, data.ConnChnlPH = AcquireKind, string(r.ConnChnlPH)
		data.ContExpCK = statesFromExpRec(key, r.ContExp, dto)
	case termdef.ExpAcceptRec:
		kind, data.ConnChnlPH = AcceptKind, string(r.ConnChnlPH)
		data.ContExpCK = statesFromExpRec(key, r.ContExp, dto)
	case termdef.ExpHireRec:
		kind, data.ConnChnlPH, data.ProcTermQN = HireKind, string(r.ConnChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
		data.ContExpCK = statesFromExpRec(key, r.ContExp, dto)
	case termdef.ExpApplyRec:
		kind, data.ConnChnlPH, data.ProcTermQN = ApplyKind, string(r.ConnChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
		data.ContExpCK = statesFromExpRec(key, r.ContExp, dto)
	case termdef.ExpReleaseRec:
		kind, data.ConnChnlPH = ReleaseKind, string(r.ConnChnlPH)
	case termdef.ExpDetachRec:
		kind, data.ConnChnlPH = DetachKind, string(r.ConnChnlPH)
	case termdef.ExpFireRec:
		kind, data.ConnChnlPH, data.ProcTermQN = FireKind, string(r.ConnChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
	case termdef.ExpQuitRec:
		kind, data.ConnChnlPH, data.ProcTermQN = QuitKind, string(r.ConnChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
	case termdef.ExpSpawnRec:
		kind, data.ProcTermQN = SpawnKind, uniqsym.ConvertToString(r.ProcTermQN)
		data.ProcCompQNs = uniqsym.ConvertToStrings(r.ProcCompQNs)
	default:
		panic(fmt.Sprintf("unexpected pool Definition record %T", rec))
	}
	dto.States = append(dto.States, ExpState{ExpCK: key, SupExpCK: parent, K: kind, Spec: data})
	return key
}

func dataToExpRec(dto ExpData) (termdef.ExpRec, error) {
	states := make(map[int64]ExpState, len(dto.States))
	for _, state := range dto.States {
		states[state.ExpCK] = state
	}
	root, ok := states[dto.ExpCK]
	if !ok {
		return nil, fmt.Errorf("pool Definition root %d is missing", dto.ExpCK)
	}
	return statesToExpRec(states, root)
}

func statesToExpRec(states map[int64]ExpState, state ExpState) (termdef.ExpRec, error) {
	key, err := contkey.ConvertFromInt(state.ExpCK)
	if err != nil {
		return nil, err
	}
	spec := state.Spec
	cont := func() (termdef.ExpRec, error) {
		child, ok := states[spec.ContExpCK]
		if !ok {
			return nil, fmt.Errorf("pool Definition child %d is missing", spec.ContExpCK)
		}
		if child.SupExpCK != state.ExpCK {
			return nil, fmt.Errorf("pool Definition child %d has unexpected parent", child.ExpCK)
		}
		return statesToExpRec(states, child)
	}
	ph := func() (symbol.ADT, error) { return symbol.ConvertFromString(spec.ConnChnlPH) }
	switch state.K {
	case AcquireKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		c, e := cont()
		return termdef.ExpAcquireRec{ExpCK: key, ConnChnlPH: p, ContExp: c}, e
	case AcceptKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		c, e := cont()
		return termdef.ExpAcceptRec{ExpCK: key, ConnChnlPH: p, ContExp: c}, e
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
		return termdef.ExpHireRec{ExpCK: key, ConnChnlPH: p, ProcTermQN: q, ContExp: c}, e
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
		return termdef.ExpApplyRec{ExpCK: key, ConnChnlPH: p, ProcTermQN: q, ContExp: c}, e
	case ReleaseKind:
		p, e := ph()
		return termdef.ExpReleaseRec{ExpCK: key, ConnChnlPH: p}, e
	case DetachKind:
		p, e := ph()
		return termdef.ExpDetachRec{ExpCK: key, ConnChnlPH: p}, e
	case FireKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		return termdef.ExpFireRec{ExpCK: key, ConnChnlPH: p, ProcTermQN: q}, e
	case QuitKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		return termdef.ExpQuitRec{ExpCK: key, ConnChnlPH: p, ProcTermQN: q}, e
	case SpawnKind:
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		if e != nil {
			return nil, e
		}
		qs, e := uniqsym.ConvertFromStrings(spec.ProcCompQNs)
		return termdef.ExpSpawnRec{ExpCK: key, ProcTermQN: q, ProcCompQNs: qs}, e
	default:
		return nil, fmt.Errorf("unexpected pool Definition kind %v", state.K)
	}
}
