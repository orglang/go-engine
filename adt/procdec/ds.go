package procdec

import (
	"orglang/go-runtime/lib/db"

	"orglang/go-runtime/adt/identity"
	"orglang/go-runtime/adt/termctx"
	"orglang/go-runtime/adt/uniqref"
)

type Repo interface {
	Insert(db.Source, DecRec) error
	SelectAll(db.Source) ([]DecRef, error)
	SelectByID(db.Source, DecRef) (DecSnap, error)
	SelectByIDs(db.Source, []identity.ADT) ([]DecRec, error)
	SelectEnv(db.Source, []identity.ADT) (map[identity.ADT]DecRec, error)
}

type decRefDS = uniqref.Data

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
