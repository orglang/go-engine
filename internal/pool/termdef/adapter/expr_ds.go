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
	ExpCK  int64      `db:"expr_ck"`
	States []ExpState `db:"states"`
}

type ExpState struct {
	ExpCK    int64        `db:"expr_ck"`
	SupExprCK int64        `db:"sup_expr_ck"`
	K        int16        `db:"kind"`
	Spec     ExpStateSpec `db:"spec"`
}

type ExpStateSpec struct {
	CommChnlPH  string   `json:"comm_chnl_ph,omitempty"`
	ProcTermQN  string   `json:"proc_term_qn,omitempty"`
	ProcCompQNs []string `json:"proc_comp_qns,omitempty"`
	ContExpCK   int64    `json:"cont_expr_ck,omitempty"`
}

func dataFromExprRec(rec termdef.ExpRec) (ExpData, error) {
	if rec == nil {
		return ExpData{}, errors.New("nil pool Definition record")
	}
	dto := ExpData{ExpCK: contkey.ConvertToInt(rec.ContentKey())}
	statesFromExprRec(0, rec, &dto)
	return dto, nil
}

func statesFromExprRec(parent int64, rec termdef.ExpRec, dto *ExpData) int64 {
	key := contkey.ConvertToInt(rec.ContentKey())
	var kind int16
	data := ExpStateSpec{}
	switch r := rec.(type) {
	case termdef.ExpAcquireRec:
		kind, data.CommChnlPH = AcquireKind, string(r.CommChnlPH)
		data.ContExpCK = statesFromExprRec(key, r.ContExp, dto)
	case termdef.ExpAcceptRec:
		kind, data.CommChnlPH = AcceptKind, string(r.CommChnlPH)
		data.ContExpCK = statesFromExprRec(key, r.ContExp, dto)
	case termdef.ExpHireRec:
		kind, data.CommChnlPH, data.ProcTermQN = HireKind, string(r.CommChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
		data.ContExpCK = statesFromExprRec(key, r.ContExp, dto)
	case termdef.ExpApplyRec:
		kind, data.CommChnlPH, data.ProcTermQN = ApplyKind, string(r.CommChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
		data.ContExpCK = statesFromExprRec(key, r.ContExp, dto)
	case termdef.ExpReleaseRec:
		kind, data.CommChnlPH = ReleaseKind, string(r.CommChnlPH)
	case termdef.ExpDetachRec:
		kind, data.CommChnlPH = DetachKind, string(r.CommChnlPH)
	case termdef.ExpFireRec:
		kind, data.CommChnlPH, data.ProcTermQN = FireKind, string(r.CommChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
	case termdef.ExpQuitRec:
		kind, data.CommChnlPH, data.ProcTermQN = QuitKind, string(r.CommChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
	case termdef.ExpSpawnRec:
		kind, data.ProcTermQN = SpawnKind, uniqsym.ConvertToString(r.ProcTermQN)
		data.ProcCompQNs = uniqsym.ConvertToStrings(r.ProcCompQNs)
	default:
		panic(fmt.Sprintf("unexpected pool Definition record %T", rec))
	}
	dto.States = append(dto.States, ExpState{ExpCK: key, SupExprCK: parent, K: kind, Spec: data})
	return key
}

func dataToExprRec(dto ExpData) (termdef.ExpRec, error) {
	states := make(map[int64]ExpState, len(dto.States))
	for _, state := range dto.States {
		if _, exists := states[state.ExpCK]; exists {
			return nil, fmt.Errorf("pool Definition contains duplicate expression %d", state.ExpCK)
		}
		states[state.ExpCK] = state
	}
	root, ok := states[dto.ExpCK]
	if !ok {
		return nil, fmt.Errorf("pool Definition root %d is missing", dto.ExpCK)
	}
	if root.SupExprCK != 0 {
		return nil, fmt.Errorf("pool Definition root %d has unexpected parent %d", root.ExpCK, root.SupExprCK)
	}
	visiting := make(map[int64]struct{}, len(states))
	visited := make(map[int64]struct{}, len(states))
	rec, err := statesToExprRec(states, root, visiting, visited)
	if err != nil {
		return nil, err
	}
	if len(visited) != len(states) {
		return nil, fmt.Errorf("pool Definition contains %d unreachable expressions", len(states)-len(visited))
	}
	return rec, nil
}

func statesToExprRec(states map[int64]ExpState, state ExpState, visiting, visited map[int64]struct{}) (termdef.ExpRec, error) {
	if _, exists := visiting[state.ExpCK]; exists {
		return nil, fmt.Errorf("pool Definition contains cyclic expression %d", state.ExpCK)
	}
	visiting[state.ExpCK] = struct{}{}
	visited[state.ExpCK] = struct{}{}
	defer delete(visiting, state.ExpCK)
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
		if child.SupExprCK != state.ExpCK {
			return nil, fmt.Errorf("pool Definition child %d has unexpected parent", child.ExpCK)
		}
		return statesToExprRec(states, child, visiting, visited)
	}
	ph := func() (symbol.ADT, error) { return symbol.ConvertFromString(spec.CommChnlPH) }
	switch state.K {
	case AcquireKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		c, e := cont()
		return termdef.ExpAcquireRec{ExpCK: key, CommChnlPH: p, ContExp: c}, e
	case AcceptKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		c, e := cont()
		return termdef.ExpAcceptRec{ExpCK: key, CommChnlPH: p, ContExp: c}, e
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
		return termdef.ExpHireRec{ExpCK: key, CommChnlPH: p, ProcTermQN: q, ContExp: c}, e
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
		return termdef.ExpApplyRec{ExpCK: key, CommChnlPH: p, ProcTermQN: q, ContExp: c}, e
	case ReleaseKind:
		p, e := ph()
		return termdef.ExpReleaseRec{ExpCK: key, CommChnlPH: p}, e
	case DetachKind:
		p, e := ph()
		return termdef.ExpDetachRec{ExpCK: key, CommChnlPH: p}, e
	case FireKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		return termdef.ExpFireRec{ExpCK: key, CommChnlPH: p, ProcTermQN: q}, e
	case QuitKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		return termdef.ExpQuitRec{ExpCK: key, CommChnlPH: p, ProcTermQN: q}, e
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
