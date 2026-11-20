package adapter

import compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"

type queryBuilder interface {
	InsertRec(string, compvaradapter.VarRecDS) (string, []any)
}

const (
	poolVars       string = "pool_comp_vars"
	poolStructVars string = "pool_struct_vars"
	poolLinearVars string = "pool_linear_vars"
)
