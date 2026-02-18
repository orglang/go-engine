package pooldec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/poolbind"
)

type repo interface {
	InsertRec(db.Source, DecRec) error
}

type decRecDS struct {
	PoolID     string               `db:"desc_id"`
	ClientBRs  []poolbind.BindRecDS `db:"client_brs"`
	ProviderBR poolbind.BindRecDS   `db:"provider_br"`
}
