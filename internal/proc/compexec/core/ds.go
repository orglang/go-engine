package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/compvar"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ExecRec) error
	ModifyRec(context.Context, tx.Conn, ExecMod) error
	GetSnapByRef(context.Context, tx.Conn, compsem.SemRef) (ExecSnap, error)
}

type ExecRecDS struct {
	CompID   string `db:"comp_id"`
	CompRN   int64  `db:"comp_rn"`
	LiabMode int16  `db:"liab_mode"`
}

type ExecModDS struct {
	CompRefs   []compsem.SemRefDS
	LinearVars []compvar.VarRecDS
}
