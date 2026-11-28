package core

import (
	"context"

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
	TermQN    uniqsym.ADT
	TermID    identity.ADT
	ProgRef   progroot.RootRef
	LiabVar   termvar.VarRec
	AssetVars []termvar.VarRec
}

type API interface {
	Create(context.Context, DecSpec) (identity.ADT, progroot.RootRef, error)
}

type Repo interface {
	AddRec(context.Context, tx.Conn, DecRec) error
	GetRecByQN(context.Context, tx.Conn, uniqsym.ADT) (DecRec, error)
	GetQNByID(context.Context, tx.Conn, identity.ADT) (uniqsym.ADT, error)
}
