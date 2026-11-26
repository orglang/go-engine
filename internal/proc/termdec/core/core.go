package core

import (
	"context"
	"iter"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/termsem"
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
	TermRef   termsem.SemRef
	TermQN    uniqsym.ADT
	LiabVar   termvar.VarRec
	AssetVars []termvar.VarRec
}

type DecSnap struct {
	TermRef   termsem.SemRef
	LiabVar   termvar.VarRec
	AssetVars []termvar.VarRec
}

type Repo interface {
	AddRec(context.Context, tx.Conn, DecRec) error
	GetRefs(context.Context, tx.Conn) ([]termsem.SemRef, error)
	GetQNByID(context.Context, tx.Conn, identity.ADT) (uniqsym.ADT, error)
	GetSnap(context.Context, tx.Conn, termsem.SemRef) (DecSnap, error)
	GetRecs(context.Context, tx.Conn, []identity.ADT) ([]DecRec, error)
	SelectEnv(context.Context, tx.Conn, []identity.ADT) (map[identity.ADT]DecRec, error)
}

func CollectEnv(recs iter.Seq[DecRec]) []uniqsym.ADT {
	return []uniqsym.ADT{}
}

// API exposes the process Declaration application operations.
type API interface {
	Incept(context.Context, uniqsym.ADT) (termsem.SemRef, error)
	Create(context.Context, DecSpec) (DecSnap, error)
	RetrieveSnap(context.Context, termsem.SemRef) (DecSnap, error)
	RetreiveRefs(context.Context) ([]termsem.SemRef, error)
}
