package pooldec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descvar"
)

type repo interface {
	InsertRec(db.Source, DecRec) error
}

type decRecDS struct {
	PoolID     string             `db:"desc_id"`
	ClientVRs  []descvar.VarRecDS `db:"client_vrs"`
	ProviderVR descvar.VarRecDS   `db:"provider_vr"`
}
