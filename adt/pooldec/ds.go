package pooldec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descvar"
)

type Repo interface {
	InsertRec(db.Source, DecRec) error
}

type decRecDS struct {
	DescID     string             `db:"desc_id"`
	DescRN     int64              `db:"desc_rn"`
	ProviderVR descvar.VarRecDS   `db:"provider_vr"`
	ClientVRs  []descvar.VarRecDS `db:"client_vrs"`
}
