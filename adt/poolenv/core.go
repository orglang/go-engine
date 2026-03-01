package poolenv

import (
	"orglang/go-engine/adt/pooldec"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
	"orglang/go-engine/adt/xactdef"
	"orglang/go-engine/adt/xactexp"
)

type EnvSpec struct {
	ProcDescQNs []uniqsym.ADT
	ProcImplQNs []uniqsym.ADT
}

type EnvSnap struct {
	XactDefs map[uniqsym.ADT]xactdef.DefSnap
	PoolDecs map[uniqsym.ADT]pooldec.DecSnap
	XactExps map[valkey.ADT]xactexp.ExpSpec
}
