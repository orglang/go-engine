package procdec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/descvar"
	"orglang/go-engine/adt/identity"
)

type Repo interface {
	InsertRec(db.Source, DecRec) error
	SelectRefs(db.Source) ([]descsem.SemRef, error)
	SelectSnap(db.Source, descsem.SemRef) (DecSnap, error)
	SelectRecs(db.Source, []identity.ADT) ([]DecRec, error)
	SelectEnv(db.Source, []identity.ADT) (map[identity.ADT]DecRec, error)
}

type decRecDS struct {
	DescID     string             `db:"desc_id"`
	DescRN     int64              `db:"desc_rn"`
	ProviderVR descvar.VarRecDS   `db:"provider_vr"`
	ClientVRs  []descvar.VarRecDS `db:"client_vrs"`
}

type decSnapDS struct {
	DescID     string             `db:"desc_id"`
	DescRN     int64              `db:"desc_rn"`
	ProviderVR descvar.VarRecDS   `db:"provider_vr"`
	ClientVRs  []descvar.VarRecDS `db:"client_vrs"`
}
