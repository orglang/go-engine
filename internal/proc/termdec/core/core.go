package core

import (
	"context"
	"iter"

	"orglang/go-engine/internal/adt/progroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
)

type DecSpec struct {
	TermQN    uniqsym.ADT
	LiabVar   termvar.VarSpec
	AssetVars []termvar.VarSpec
}

type DecRec struct {
	TermID    identity.ADT
	ProgRef   progroot.RootRef
	TermQN    uniqsym.ADT
	LiabVar   termvar.VarRec
	AssetVars []termvar.VarRec
}

type DecSnap struct {
	TermID    identity.ADT
	ProgRef   progroot.RootRef
	LiabVar   termvar.VarRec
	AssetVars []termvar.VarRec
}

type Repo interface {
	AddRec(context.Context, tx.Conn, DecRec) error
	GetIDs(context.Context, tx.Conn) ([]identity.ADT, error)
	GetQNByID(context.Context, tx.Conn, identity.ADT) (uniqsym.ADT, error)
	GetSnap(context.Context, tx.Conn, identity.ADT) (DecSnap, error)
	GetRecs(context.Context, tx.Conn, []identity.ADT) ([]DecRec, error)
	SelectEnv(context.Context, tx.Conn, []identity.ADT) (map[identity.ADT]DecRec, error)
}

func CollectEnv(recs iter.Seq[DecRec]) []uniqsym.ADT { return []uniqsym.ADT{} }

type API interface {
	Incept(context.Context, uniqsym.ADT) (identity.ADT, progroot.RootRef, error)
	Create(context.Context, DecSpec) (DecSnap, error)
	RetrieveSnap(context.Context, identity.ADT) (DecSnap, error)
	RetreiveRefs(context.Context) ([]identity.ADT, error)
}
