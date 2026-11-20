package adapter

import adtcompvar "orglang/go-engine/internal/adt/compvar"

type queryBuilder interface {
	InsertRec(string, adtcompvar.VarRecDS) (string, []any)
}

const (
	poolVars       string = "pool_comp_vars"
	poolStructVars string = "pool_struct_vars"
	poolLinearVars string = "pool_linear_vars"
)
