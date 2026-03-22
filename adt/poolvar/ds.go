package poolvar

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/implvar"
)

type Repo interface {
	AddRecs(db.Source, []implvar.VarRec) error
}
