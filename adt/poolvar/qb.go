package poolvar

import (
	"orglang/go-engine/adt/implvar"
)

const (
	poolVars       string = "pool_vars"
	poolStructVars string = "pool_struct_vars"
	poolLinearVars string = "pool_linear_vars"
)

type queryBuilder interface {
	insertRec(string, implvar.VarRecDS) (string, []any)
}
