package procdec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/procbind"
	"orglang/go-engine/adt/uniqref"
)

type Repo interface {
	InsertRec(db.Source, DecRec) error
	SelectRefs(db.Source) ([]DecRef, error)
	SelectSnap(db.Source, DecRef) (DecSnap, error)
	SelectRecs(db.Source, []identity.ADT) ([]DecRec, error)
	SelectEnv(db.Source, []identity.ADT) (map[identity.ADT]DecRec, error)
}

type decRefDS = uniqref.Data

type decRecDS struct {
	ID         string                `db:"dec_id"`
	RN         int64                 `db:"dec_rn"`
	SynVK      int64                 `db:"syn_vk"`
	ClientBSes []procbind.BindSpecDS `db:"client_bses"`
	ProviderBS procbind.BindSpecDS   `db:"provider_bs"`
}

type decSnapDS struct {
	ID         string                `db:"dec_id"`
	RN         int64                 `db:"dec_rn"`
	ClientBSes []procbind.BindSpecDS `db:"client_bses"`
	ProviderBS procbind.BindSpecDS   `db:"provider_bs"`
}
