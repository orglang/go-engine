package compvar

import (
	"orglang/go-engine/adt/compvar"
)

const (
	poolVars       string = "pool_vars"
	poolStructVars string = "pool_struct_vars"
	poolLinearVars string = "pool_linear_vars"
)

type queryBuilder interface {
	insertRec(string, compvar.VarRecDS) (string, []any)
}
