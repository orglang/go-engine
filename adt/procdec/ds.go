package procdec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implvar"
)

type Repo interface {
	InsertRec(db.Source, DecRec) error
	SelectRefs(db.Source) ([]descsem.SemRef, error)
	SelectSnap(db.Source, descsem.SemRef) (DecSnap, error)
	SelectRecs(db.Source, []identity.ADT) ([]DecRec, error)
	SelectEnv(db.Source, []identity.ADT) (map[identity.ADT]DecRec, error)
}

type decRecDS struct {
	DescID     string              `db:"desc_id"`
	DescRN     int64               `db:"desc_rn"`
	ClientVSes []implvar.VarSpecDS `db:"client_vrs"`
	ProviderVS implvar.VarSpecDS   `db:"provider_vr"`
}

type decSnapDS struct {
	DescID     string              `db:"desc_id"`
	DescRN     int64               `db:"desc_rn"`
	ClientVSes []implvar.VarSpecDS `db:"client_vrs"`
	ProviderVS implvar.VarSpecDS   `db:"provider_vr"`
}
