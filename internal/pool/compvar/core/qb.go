package core

import (
	"orglang/go-engine/adt/compvar"
)

type QueryBuilder interface {
	InsertRec(string, compvar.VarRecDS) (string, []any)
}
