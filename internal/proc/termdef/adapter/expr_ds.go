package adapter

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	termdef "orglang/go-engine/internal/proc/termdef/core"
)

const (
	CloseKind int16 = iota + 1
	WaitKind
	SendKind
	RecvKind
	LabKind
	CaseKind
	LinkKind
	FwdKind
	CallKind
	SpawnKind
	AcquireKind
	AcceptKind
	DetachKind
	ReleaseKind
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
	ContChnlPH string           `json:"cont_chnl_ph,omitempty"`
	CommChnlPH string           `json:"comm_chnl_ph,omitempty"`
	ValChnlPH  string           `json:"val_chnl_ph,omitempty"`
	ValLabQN   string           `json:"val_lab_qn,omitempty"`
	NewChnlPH  string           `json:"new_chnl_ph,omitempty"`
	ContExpCK  int64            `json:"cont_expr_ck,omitempty"`
	ContExpCKs map[string]int64 `json:"cont_expr_cks,omitempty"`
	ProcTermQN string           `json:"proc_term_qn,omitempty"`
	X          string           `json:"x,omitempty"`
	Ys         []string         `json:"ys,omitempty"`
	ValChnlPHs []string         `json:"val_chnl_phs,omitempty"`
	NewChnlPHs []string         `json:"new_chnl_phs,omitempty"`
}

func dataFromExprRec(rec termdef.ExpRec) (ExpData, error) {
	if rec == nil {
		return ExpData{}, errors.New("nil process Definition record")
	}
	dto := ExpData{ExpCK: contkey.ConvertToInt(rec.ContentKey())}
	statesFromExprRec(0, rec, &dto)
	return dto, nil
}
func statesFromExprRec(parent int64, rec termdef.ExpRec, dto *ExpData) int64 {
	key := contkey.ConvertToInt(rec.ContentKey())
	var kind int16
	d := ExpStateSpec{}
	child := func(r termdef.ExpRec) int64 { return statesFromExprRec(key, r, dto) }
	switch r := rec.(type) {
	case termdef.ExpCloseRec:
		kind = CloseKind
		d.ContChnlPH = string(r.ContChnlPH)
	case termdef.ExpWaitRec:
		kind = WaitKind
		d.ContChnlPH = string(r.ContChnlPH)
		d.ContExpCK = child(r.ContExp)
	case termdef.ExpSendRec:
		kind = SendKind
		d.CommChnlPH = string(r.CommChnlPH)
		d.ValChnlPH = string(r.ValChnlPH)
	case termdef.ExpRecvRec:
		kind = RecvKind
		d.CommChnlPH = string(r.CommChnlPH)
		d.NewChnlPH = string(r.NewChnlPH)
		d.ContExpCK = child(r.ContExp)
	case termdef.ExpLabRec:
		kind = LabKind
		d.CommChnlPH = string(r.CommChnlPH)
		d.ValLabQN = uniqsym.ConvertToString(r.ValLabQN)
		d.ContExpCK = child(r.ContExp)
	case termdef.ExpCaseRec:
		kind = CaseKind
		d.CommChnlPH = string(r.CommChnlPH)
		d.ContExpCKs = make(map[string]int64, len(r.ContExprs))
		labels := make([]uniqsym.ADT, 0, len(r.ContExprs))
		for q := range r.ContExprs {
			labels = append(labels, q)
		}
		slices.SortFunc(labels, func(a, b uniqsym.ADT) int {
			return strings.Compare(uniqsym.ConvertToString(a), uniqsym.ConvertToString(b))
		})
		for _, q := range labels {
			d.ContExpCKs[uniqsym.ConvertToString(q)] = child(r.ContExprs[q])
		}
	case termdef.ExpLinkRec:
		kind = LinkKind
		d.ProcTermQN = uniqsym.ConvertToString(r.ProcTermQN)
		d.X = identity.ConvertToString(r.X)
		d.Ys = make([]string, len(r.Ys))
		for i, v := range r.Ys {
			d.Ys[i] = identity.ConvertToString(v)
		}
	case termdef.ExpFwdRec:
		kind = FwdKind
		d.CommChnlPH = string(r.CommChnlPH)
		d.ContChnlPH = string(r.ContChnlPH)
	case termdef.ExpCallRec:
		kind = CallKind
		d.CommChnlPH = string(r.CommChnlPH)
		d.NewChnlPH = string(r.NewChnlPH)
		d.ProcTermQN = uniqsym.ConvertToString(r.ProcTermQN)
		d.ValChnlPHs = make([]string, len(r.ValChnlPHs))
		for i, v := range r.ValChnlPHs {
			d.ValChnlPHs[i] = string(v)
		}
		d.ContExpCK = child(r.ContExp)
	case termdef.ExpSpawnRec:
		kind = SpawnKind
		d.CommChnlPH = string(r.CommChnlPH)
		d.ProcTermQN = uniqsym.ConvertToString(r.ProcTermQN)
		d.NewChnlPHs = make([]string, len(r.NewChnlPHs))
		for i, v := range r.NewChnlPHs {
			d.NewChnlPHs[i] = string(v)
		}
		d.ContExpCK = child(r.ContExp)
	case termdef.ExpAcquireRec:
		kind = AcquireKind
		d.CommChnlPH = string(r.CommChnlPH)
		d.ContExpCK = child(r.ContExp)
	case termdef.ExpAcceptRec:
		kind = AcceptKind
		d.CommChnlPH = string(r.CommChnlPH)
		d.ContExpCK = child(r.ContExp)
	case termdef.ExpDetachRec:
		kind = DetachKind
		d.CommChnlPH = string(r.CommChnlPH)
	case termdef.ExpReleaseRec:
		kind = ReleaseKind
		d.CommChnlPH = string(r.CommChnlPH)
	default:
		panic(fmt.Sprintf("unexpected process Definition record %T", rec))
	}
	dto.States = append(dto.States, ExpState{ExpCK: key, SupExprCK: parent, K: kind, Spec: d})
	return key
}

func dataToExprRec(dto ExpData) (termdef.ExpRec, error) {
	states := make(map[int64]ExpState, len(dto.States))
	for _, state := range dto.States {
		if _, exists := states[state.ExpCK]; exists {
			return nil, fmt.Errorf("process Definition contains duplicate expression %d", state.ExpCK)
		}
		states[state.ExpCK] = state
	}
	root, ok := states[dto.ExpCK]
	if !ok {
		return nil, fmt.Errorf("process Definition root %d is missing", dto.ExpCK)
	}
	if root.SupExprCK != 0 {
		return nil, fmt.Errorf("process Definition root %d has unexpected parent %d", root.ExpCK, root.SupExprCK)
	}
	visiting := make(map[int64]struct{}, len(states))
	visited := make(map[int64]struct{}, len(states))
	rec, err := statesToExprRec(states, root, visiting, visited)
	if err != nil {
		return nil, err
	}
	if len(visited) != len(states) {
		return nil, fmt.Errorf("process Definition contains %d unreachable expressions", len(states)-len(visited))
	}
	return rec, nil
}
func statesToExprRec(states map[int64]ExpState, st ExpState, visiting, visited map[int64]struct{}) (termdef.ExpRec, error) {
	if _, exists := visiting[st.ExpCK]; exists {
		return nil, fmt.Errorf("process Definition contains cyclic expression %d", st.ExpCK)
	}
	visiting[st.ExpCK] = struct{}{}
	visited[st.ExpCK] = struct{}{}
	defer delete(visiting, st.ExpCK)
	key, err := contkey.ConvertFromInt(st.ExpCK)
	if err != nil {
		return nil, err
	}
	d := st.Spec
	child := func(ck int64) (termdef.ExpRec, error) {
		s, ok := states[ck]
		if !ok {
			return nil, fmt.Errorf("process Definition child %d is missing", ck)
		}
		if s.SupExprCK != st.ExpCK {
			return nil, fmt.Errorf("process Definition child %d has unexpected parent", ck)
		}
		return statesToExprRec(states, s, visiting, visited)
	}
	ph := symbol.ConvertFromString
	switch st.K {
	case CloseKind:
		p, e := ph(d.ContChnlPH)
		return termdef.ExpCloseRec{ExpCK: key, ContChnlPH: p}, e
	case WaitKind:
		p, e := ph(d.ContChnlPH)
		if e != nil {
			return nil, e
		}
		c, e := child(d.ContExpCK)
		return termdef.ExpWaitRec{ExpCK: key, ContChnlPH: p, ContExp: c}, e
	case SendKind:
		x, e := ph(d.CommChnlPH)
		if e != nil {
			return nil, e
		}
		y, e := ph(d.ValChnlPH)
		return termdef.ExpSendRec{ExpCK: key, CommChnlPH: x, ValChnlPH: y}, e
	case RecvKind:
		x, e := ph(d.CommChnlPH)
		if e != nil {
			return nil, e
		}
		y, e := ph(d.NewChnlPH)
		if e != nil {
			return nil, e
		}
		c, e := child(d.ContExpCK)
		return termdef.ExpRecvRec{ExpCK: key, CommChnlPH: x, NewChnlPH: y, ContExp: c}, e
	case LabKind:
		x, e := ph(d.CommChnlPH)
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(d.ValLabQN)
		if e != nil {
			return nil, e
		}
		c, e := child(d.ContExpCK)
		return termdef.ExpLabRec{ExpCK: key, CommChnlPH: x, ValLabQN: q, ContExp: c}, e
	case CaseKind:
		x, e := ph(d.CommChnlPH)
		if e != nil {
			return nil, e
		}
		cs := make(map[uniqsym.ADT]termdef.ExpRec, len(d.ContExpCKs))
		for name, ck := range d.ContExpCKs {
			q, e := uniqsym.ConvertFromString(name)
			if e != nil {
				return nil, e
			}
			c, e := child(ck)
			if e != nil {
				return nil, e
			}
			cs[q] = c
		}
		return termdef.ExpCaseRec{ExpCK: key, CommChnlPH: x, ContExprs: cs}, nil
	case LinkKind:
		q, e := uniqsym.ConvertFromString(d.ProcTermQN)
		if e != nil {
			return nil, e
		}
		x, e := identity.ConvertFromString(d.X)
		if e != nil {
			return nil, e
		}
		ys := make([]identity.ADT, len(d.Ys))
		for i, v := range d.Ys {
			ys[i], e = identity.ConvertFromString(v)
			if e != nil {
				return nil, e
			}
		}
		return termdef.ExpLinkRec{ExpCK: key, ProcTermQN: q, X: x, Ys: ys}, nil
	case FwdKind:
		x, e := ph(d.CommChnlPH)
		if e != nil {
			return nil, e
		}
		y, e := ph(d.ContChnlPH)
		return termdef.ExpFwdRec{ExpCK: key, CommChnlPH: x, ContChnlPH: y}, e
	case CallKind:
		x, e := ph(d.CommChnlPH)
		if e != nil {
			return nil, e
		}
		y, e := ph(d.NewChnlPH)
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(d.ProcTermQN)
		if e != nil {
			return nil, e
		}
		chs := make([]symbol.ADT, len(d.ValChnlPHs))
		for i, v := range d.ValChnlPHs {
			chs[i], e = ph(v)
			if e != nil {
				return nil, e
			}
		}
		c, e := child(d.ContExpCK)
		return termdef.ExpCallRec{ExpCK: key, CommChnlPH: x, NewChnlPH: y, ProcTermQN: q, ValChnlPHs: chs, ContExp: c}, e
	case SpawnKind:
		x, e := ph(d.CommChnlPH)
		if e != nil {
			return nil, e
		}
		q, e := uniqsym.ConvertFromString(d.ProcTermQN)
		if e != nil {
			return nil, e
		}
		chs := make([]symbol.ADT, len(d.NewChnlPHs))
		for i, v := range d.NewChnlPHs {
			chs[i], e = ph(v)
			if e != nil {
				return nil, e
			}
		}
		c, e := child(d.ContExpCK)
		return termdef.ExpSpawnRec{ExpCK: key, CommChnlPH: x, ProcTermQN: q, NewChnlPHs: chs, ContExp: c}, e
	case AcquireKind:
		x, e := ph(d.CommChnlPH)
		if e != nil {
			return nil, e
		}
		c, e := child(d.ContExpCK)
		return termdef.ExpAcquireRec{ExpCK: key, CommChnlPH: x, ContExp: c}, e
	case AcceptKind:
		x, e := ph(d.CommChnlPH)
		if e != nil {
			return nil, e
		}
		c, e := child(d.ContExpCK)
		return termdef.ExpAcceptRec{ExpCK: key, CommChnlPH: x, ContExp: c}, e
	case DetachKind:
		x, e := ph(d.CommChnlPH)
		return termdef.ExpDetachRec{ExpCK: key, CommChnlPH: x}, e
	case ReleaseKind:
		x, e := ph(d.CommChnlPH)
		return termdef.ExpReleaseRec{ExpCK: key, CommChnlPH: x}, e
	default:
		return nil, fmt.Errorf("unexpected process Definition kind %d", st.K)
	}
}
