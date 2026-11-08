package core

import (
	"orglang/go-engine/internal/adt/compvar"
)

type QueryBuilder interface {
	InsertRec(string, compvar.VarRecDS) (string, []any)
}
