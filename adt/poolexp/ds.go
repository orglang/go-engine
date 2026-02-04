package poolexp

import (
	"orglang/go-engine/lib/db"
)

type Repo interface {
	InsertRec(db.Source, ExpSpec) error
}
