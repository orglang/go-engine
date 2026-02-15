package syndec

import (
	"orglang/go-engine/lib/db"
)

type Repo interface {
	InsertRec(db.Source, DecRec) error
}

type decRecDS struct {
	DecID string
	DecRN int64
	DecQN string
}
