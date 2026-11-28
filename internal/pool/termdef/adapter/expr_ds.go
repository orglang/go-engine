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

type ExprData struct {
	ExprCK  int64      `db:"expr_ck"`
	States []ExprState `db:"states"`
}

type ExprState struct {
	ExprCK     int64        `db:"expr_ck"`
	SupExprCK int64        `db:"sup_expr_ck"`
	K         int16        `db:"kind"`
	Spec      ExprStateSpec `db:"spec"`
}

type ExprStateSpec struct {
	CommChnlPH  string   `json:"comm_chnl_ph,omitempty"`
	ProcTermQN  string   `json:"proc_term_qn,omitempty"`
	ProcCompQNs []string `json:"proc_comp_qns,omitempty"`
	ContExpCK   int64    `json:"cont_expr_ck,omitempty"`
}

func dataFromExprRec(rec termdef.ExprRec) (ExprData, error) {
	if rec == nil {
		return ExprData{}, errors.New("nil pool Definition record")
	}
	dto := ExprData{ExprCK: contkey.ConvertToInt(rec.ContentKey())}
	statesFromExprRec(0, rec, &dto)
	return dto, nil
}

func statesFromExprRec(parent int64, rec termdef.ExprRec, dto *ExprData) int64 {
	key := contkey.ConvertToInt(rec.ContentKey())
	var kind int16
	data := ExprStateSpec{}
	switch r := rec.(type) {
	case termdef.ExprAcquireRec:
		kind, data.CommChnlPH = AcquireKind, string(r.CommChnlPH)
		data.ContExpCK = statesFromExprRec(key, r.ContExp, dto)
	case termdef.ExprAcceptRec:
		kind, data.CommChnlPH = AcceptKind, string(r.CommChnlPH)
		data.ContExpCK = statesFromExprRec(key, r.ContExp, dto)
	case termdef.ExprHireRec:
		kind, data.CommChnlPH, data.ProcTermQN = HireKind, string(r.CommChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
		data.ContExpCK = statesFromExprRec(key, r.ContExp, dto)
	case termdef.ExprApplyRec:
		kind, data.CommChnlPH, data.ProcTermQN = ApplyKind, string(r.CommChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
		data.ContExpCK = statesFromExprRec(key, r.ContExp, dto)
	case termdef.ExprReleaseRec:
		kind, data.CommChnlPH = ReleaseKind, string(r.CommChnlPH)
	case termdef.ExprDetachRec:
		kind, data.CommChnlPH = DetachKind, string(r.CommChnlPH)
	case termdef.ExprFireRec:
		kind, data.CommChnlPH, data.ProcTermQN = FireKind, string(r.CommChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
	case termdef.ExprQuitRec:
		kind, data.CommChnlPH, data.ProcTermQN = QuitKind, string(r.CommChnlPH), uniqsym.ConvertToString(r.ProcTermQN)
	case termdef.ExprSpawnRec:
		kind, data.ProcTermQN = SpawnKind, uniqsym.ConvertToString(r.ProcTermQN)
		data.ProcCompQNs = uniqsym.ConvertToStrings(r.ProcCompQNs)
	default:
		panic(fmt.Sprintf("unexpected pool Definition record %T", rec))
	}
	dto.States = append(dto.States, ExprState{ExprCK: key, SupExprCK: parent, K: kind, Spec: data})
	return key
}

func dataToExprRec(dto ExprData) (termdef.ExprRec, error) {
	states := make(map[int64]ExprState, len(dto.States))
	for _, state := range dto.States {
		if _, exists := states[state.ExprCK]; exists {
			return nil, fmt.Errorf("pool Definition contains duplicate expression %d", state.ExprCK)
		}
		states[state.ExprCK] = state
	}
	root, ok := states[dto.ExprCK]
	if !ok {
		return nil, fmt.Errorf("pool Definition root %d is missing", dto.ExprCK)
	}
	if root.SupExprCK != 0 {
		return nil, fmt.Errorf("pool Definition root %d has unexpected parent %d", root.ExprCK, root.SupExprCK)
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

func statesToExprRec(states map[int64]ExprState, state ExprState, visiting, visited map[int64]struct{}) (termdef.ExprRec, error) {
	if _, exists := visiting[state.ExprCK]; exists {
		return nil, fmt.Errorf("pool Definition contains cyclic expression %d", state.ExprCK)
	}
	visiting[state.ExprCK] = struct{}{}
	visited[state.ExprCK] = struct{}{}
	defer delete(visiting, state.ExprCK)
	key, err := contkey.ConvertFromInt(state.ExprCK)
	if err != nil {
		return nil, err
	}
	spec := state.Spec
	cont := func() (termdef.ExprRec, error) {
		child, ok := states[spec.ContExpCK]
		if !ok {
			return nil, fmt.Errorf("pool Definition child %d is missing", spec.ContExpCK)
		}
		if child.SupExprCK != state.ExprCK {
			return nil, fmt.Errorf("pool Definition child %d has unexpected parent", child.ExprCK)
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
		return termdef.ExprAcquireRec{ExprCK: key, CommChnlPH: p, ContExp: c}, e
	case AcceptKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		c, e := cont()
		return termdef.ExprAcceptRec{ExprCK: key, CommChnlPH: p, ContExp: c}, e
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
		return termdef.ExprHireRec{ExprCK: key, CommChnlPH: p, ProcTermQN: q, ContExp: c}, e
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
		return termdef.ExprApplyRec{ExprCK: key, CommChnlPH: p, ProcTermQN: q, ContExp: c}, e
	case ReleaseKind:
		p, e := ph()
		return termdef.ExprReleaseRec{ExprCK: key, CommChnlPH: p}, e
	case DetachKind:
		p, e := ph()
		return termdef.ExprDetachRec{ExprCK: key, CommChnlPH: p}, e
	case FireKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		return termdef.ExprFireRec{ExprCK: key, CommChnlPH: p, ProcTermQN: q}, e
	case QuitKind:
		p, e := ph()
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		return termdef.ExprQuitRec{ExprCK: key, CommChnlPH: p, ProcTermQN: q}, e
	case SpawnKind:
		q, e := uniqsym.ConvertFromString(spec.ProcTermQN)
		if e != nil {
			return nil, e
		}
		qs, e := uniqsym.ConvertFromStrings(spec.ProcCompQNs)
		return termdef.ExprSpawnRec{ExprCK: key, ProcTermQN: q, ProcCompQNs: qs}, e
	default:
		return nil, fmt.Errorf("unexpected pool Definition kind %v", state.K)
	}
}
