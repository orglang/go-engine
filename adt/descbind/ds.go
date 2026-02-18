package descbind

import (
	"orglang/go-engine/lib/db"
)

type Repo interface {
	InsertRec(db.Source, BindRec) error
}

type bindRecDS struct {
	DescQN string `db:"desc_qn"`
	DescID string `db:"desc_id"`
}
