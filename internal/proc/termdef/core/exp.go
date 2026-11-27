package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

// ExpSpec is one node in a Term Definition tree.
type ExpSpec interface {
	expSpec()
}

type ExpCloseSpec struct{ ContChnlPH symbol.ADT }

func (ExpCloseSpec) expSpec() {}

type ExpWaitSpec struct {
	ContChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (ExpWaitSpec) expSpec() {}

type ExpSendSpec struct {
	ConnChnlPH symbol.ADT
	ValChnlPH  symbol.ADT
}

func (ExpSendSpec) expSpec() {}

type ExpRecvSpec struct {
	ConnChnlPH symbol.ADT
	NewChnlPH  symbol.ADT
	ContExp    ExpSpec
}

func (ExpRecvSpec) expSpec() {}

type ExpLabSpec struct {
	ConnChnlPH symbol.ADT
	ValLabQN   uniqsym.ADT
	ContExp    ExpSpec
}

func (ExpLabSpec) expSpec() {}

type ExpCaseSpec struct {
	ConnChnlPH symbol.ADT
	ContExps   map[uniqsym.ADT]ExpSpec
}

func (ExpCaseSpec) expSpec() {}

type ExpLinkSpec struct {
	ProcTermQN uniqsym.ADT
	X          identity.ADT
	Ys         []identity.ADT
}

func (ExpLinkSpec) expSpec() {}

type ExpFwdSpec struct {
	ConnChnlPH symbol.ADT
	ContChnlPH symbol.ADT
}

func (ExpFwdSpec) expSpec() {}

type ExpCallSpec struct {
	ConnChnlPH symbol.ADT
	NewChnlPH  symbol.ADT
	ProcTermQN uniqsym.ADT
	ValChnlPHs []symbol.ADT
	ContExp    ExpSpec
}

func (ExpCallSpec) expSpec() {}

type ExpSpawnSpec struct {
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	NewChnlPHs []symbol.ADT
	ContExp    ExpSpec
}

func (ExpSpawnSpec) expSpec() {}

type ExpAcquireSpec struct {
	ConnChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (ExpAcquireSpec) expSpec() {}

type ExpAcceptSpec struct {
	ConnChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (ExpAcceptSpec) expSpec() {}

type ExpDetachSpec struct{ ConnChnlPH symbol.ADT }

func (ExpDetachSpec) expSpec() {}

type ExpReleaseSpec struct{ ConnChnlPH symbol.ADT }

func (ExpReleaseSpec) expSpec() {}

// ExpRec is the recursive, content-addressed domain representation of a Definition.
type ExpRec interface {
	expRec()
	ContentKey() contkey.ADT
}
type ExpCloseRec struct {
	ExpCK      contkey.ADT
	ContChnlPH symbol.ADT
}

func (ExpCloseRec) expRec()                   {}
func (r ExpCloseRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpWaitRec struct {
	ExpCK      contkey.ADT
	ContChnlPH symbol.ADT
	ContExp    ExpRec
}

func (ExpWaitRec) expRec()                   {}
func (r ExpWaitRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpSendRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ValChnlPH  symbol.ADT
}

func (ExpSendRec) expRec()                   {}
func (r ExpSendRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpRecvRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	NewChnlPH  symbol.ADT
	ContExp    ExpRec
}

func (ExpRecvRec) expRec()                   {}
func (r ExpRecvRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpLabRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ValLabQN   uniqsym.ADT
	ContExp    ExpRec
}

func (ExpLabRec) expRec()                   {}
func (r ExpLabRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpCaseRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ContExps   map[uniqsym.ADT]ExpRec
}

func (ExpCaseRec) expRec()                   {}
func (r ExpCaseRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpLinkRec struct {
	ExpCK      contkey.ADT
	ProcTermQN uniqsym.ADT
	X          identity.ADT
	Ys         []identity.ADT
}

func (ExpLinkRec) expRec()                   {}
func (r ExpLinkRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpFwdRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ContChnlPH symbol.ADT
}

func (ExpFwdRec) expRec()                   {}
func (r ExpFwdRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpCallRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	NewChnlPH  symbol.ADT
	ProcTermQN uniqsym.ADT
	ValChnlPHs []symbol.ADT
	ContExp    ExpRec
}

func (ExpCallRec) expRec()                   {}
func (r ExpCallRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpSpawnRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	NewChnlPHs []symbol.ADT
	ContExp    ExpRec
}

func (ExpSpawnRec) expRec()                   {}
func (r ExpSpawnRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpAcquireRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ContExp    ExpRec
}

func (ExpAcquireRec) expRec()                   {}
func (r ExpAcquireRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpAcceptRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ContExp    ExpRec
}

func (ExpAcceptRec) expRec()                   {}
func (r ExpAcceptRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpDetachRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
}

func (ExpDetachRec) expRec()                   {}
func (r ExpDetachRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpReleaseRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
}

func (ExpReleaseRec) expRec()                   {}
func (r ExpReleaseRec) ContentKey() contkey.ADT { return r.ExpCK }

// ConvertSpecToRec converts a declarative Definition into its content-addressed domain record.
func ConvertSpecToRec(spec ExpSpec) (ExpRec, error) {
	key, _, err := ExpContentKey(spec)
	if err != nil {
		return nil, err
	}
	switch s := spec.(type) {
	case ExpCloseSpec:
		return ExpCloseRec{key, s.ContChnlPH}, nil
	case ExpWaitSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpWaitRec{key, s.ContChnlPH, c}, nil
	case ExpSendSpec:
		return ExpSendRec{key, s.ConnChnlPH, s.ValChnlPH}, nil
	case ExpRecvSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpRecvRec{key, s.ConnChnlPH, s.NewChnlPH, c}, nil
	case ExpLabSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpLabRec{key, s.ConnChnlPH, s.ValLabQN, c}, nil
	case ExpCaseSpec:
		cs := make(map[uniqsym.ADT]ExpRec, len(s.ContExps))
		for q, v := range s.ContExps {
			c, e := ConvertSpecToRec(v)
			if e != nil {
				return nil, e
			}
			cs[q] = c
		}
		return ExpCaseRec{key, s.ConnChnlPH, cs}, nil
	case ExpLinkSpec:
		return ExpLinkRec{key, s.ProcTermQN, s.X, s.Ys}, nil
	case ExpFwdSpec:
		return ExpFwdRec{key, s.ConnChnlPH, s.ContChnlPH}, nil
	case ExpCallSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpCallRec{key, s.ConnChnlPH, s.NewChnlPH, s.ProcTermQN, s.ValChnlPHs, c}, nil
	case ExpSpawnSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpSpawnRec{key, s.ConnChnlPH, s.ProcTermQN, s.NewChnlPHs, c}, nil
	case ExpAcquireSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpAcquireRec{key, s.ConnChnlPH, c}, nil
	case ExpAcceptSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpAcceptRec{key, s.ConnChnlPH, c}, nil
	case ExpDetachSpec:
		return ExpDetachRec{key, s.ConnChnlPH}, nil
	case ExpReleaseSpec:
		return ExpReleaseRec{key, s.ConnChnlPH}, nil
	default:
		return nil, fmt.Errorf("unexpected process Definition expression spec %T", spec)
	}
}

func ConvertRecToSpec(rec ExpRec) ExpSpec {
	switch r := rec.(type) {
	case ExpCloseRec:
		return ExpCloseSpec{r.ContChnlPH}
	case ExpWaitRec:
		return ExpWaitSpec{r.ContChnlPH, ConvertRecToSpec(r.ContExp)}
	case ExpSendRec:
		return ExpSendSpec{r.ConnChnlPH, r.ValChnlPH}
	case ExpRecvRec:
		return ExpRecvSpec{r.ConnChnlPH, r.NewChnlPH, ConvertRecToSpec(r.ContExp)}
	case ExpLabRec:
		return ExpLabSpec{r.ConnChnlPH, r.ValLabQN, ConvertRecToSpec(r.ContExp)}
	case ExpCaseRec:
		cs := make(map[uniqsym.ADT]ExpSpec, len(r.ContExps))
		for q, v := range r.ContExps {
			cs[q] = ConvertRecToSpec(v)
		}
		return ExpCaseSpec{r.ConnChnlPH, cs}
	case ExpLinkRec:
		return ExpLinkSpec{r.ProcTermQN, r.X, r.Ys}
	case ExpFwdRec:
		return ExpFwdSpec{r.ConnChnlPH, r.ContChnlPH}
	case ExpCallRec:
		return ExpCallSpec{r.ConnChnlPH, r.NewChnlPH, r.ProcTermQN, r.ValChnlPHs, ConvertRecToSpec(r.ContExp)}
	case ExpSpawnRec:
		return ExpSpawnSpec{r.ConnChnlPH, r.ProcTermQN, r.NewChnlPHs, ConvertRecToSpec(r.ContExp)}
	case ExpAcquireRec:
		return ExpAcquireSpec{r.ConnChnlPH, ConvertRecToSpec(r.ContExp)}
	case ExpAcceptRec:
		return ExpAcceptSpec{r.ConnChnlPH, ConvertRecToSpec(r.ContExp)}
	case ExpDetachRec:
		return ExpDetachSpec{r.ConnChnlPH}
	case ExpReleaseRec:
		return ExpReleaseSpec{r.ConnChnlPH}
	default:
		panic(fmt.Sprintf("unexpected process Definition expression rec %T", rec))
	}
}
