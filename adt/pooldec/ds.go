package pooldec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/poolbind"
)

type repo interface {
	InsertRec(db.Source, DecRec) error
}

type decRecDS struct {
	ID         string               `db:"dec_id"`
	RN         int64                `db:"dec_rn"`
	SynVK      int64                `db:"syn_vk"`
	ClientBRs  []poolbind.BindRecDS `db:"client_brs"`
	ProviderBR poolbind.BindRecDS   `db:"provider_br"`
}
