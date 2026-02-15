package synonym

import (
	"orglang/go-engine/lib/db"
)

type Repo interface {
	InsertRec(db.Source, Rec) error
}

type recDS struct {
	SynQN string `db:"syn_qn"`
	SynVK int64  `db:"syn_vk"`
}
