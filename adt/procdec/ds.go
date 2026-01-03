package procdec

import (
	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/termctx"
)

type Repo interface {
	Insert(sd.Source, DecRec) error
	SelectAll(sd.Source) ([]DecRef, error)
	SelectByID(sd.Source, identity.ADT) (DecSnap, error)
	SelectByIDs(sd.Source, []identity.ADT) ([]DecRec, error)
	SelectEnv(sd.Source, []identity.ADT) (map[identity.ADT]DecRec, error)
}

type decRefDS struct {
	DecID string `db:"dec_id"`
	DecRN int64  `db:"dec_rn"`
}

type decRecDS struct {
	DecID string                `db:"dec_id"`
	Title string                `db:"title"`
	Ys    []termctx.BindClaimDS `db:"ys"`
	X     termctx.BindClaimDS   `db:"x"`
	DecRN int64                 `db:"dec_rn"`
}

type decSnapDS struct {
	DecID string                `db:"dec_id"`
	Title string                `db:"title"`
	Ys    []termctx.BindClaimDS `db:"ys"`
	X     termctx.BindClaimDS   `db:"x"`
	DecRN int64                 `db:"dec_rn"`
}
