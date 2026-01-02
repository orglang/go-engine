package procdec

import (
	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/termctx"
)

type Repo interface {
	Insert(sd.Source, ProcRec) error
	SelectAll(sd.Source) ([]ProcRef, error)
	SelectByID(sd.Source, identity.ADT) (ProcSnap, error)
	SelectByIDs(sd.Source, []identity.ADT) ([]ProcRec, error)
	SelectEnv(sd.Source, []identity.ADT) (map[identity.ADT]ProcRec, error)
}

type sigRefDS struct {
	SigID string `db:"sig_id"`
	SigRN int64  `db:"rev"`
	Title string `db:"title"`
}

type sigRecDS struct {
	SigID string                `db:"sig_id"`
	Title string                `db:"title"`
	Ys    []termctx.BindClaimDS `db:"ys"`
	X     termctx.BindClaimDS   `db:"x"`
	SigRN int64                 `db:"rev"`
}

type sigSnapDS struct {
	SigID string                `db:"sig_id"`
	Title string                `db:"title"`
	Ys    []termctx.BindClaimDS `db:"ys"`
	X     termctx.BindClaimDS   `db:"x"`
	SigRN int64                 `db:"rev"`
}
