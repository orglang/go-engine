package poolimpl

import (
	"orglang/go-engine/adt/implsem"
)

const (
	implSems       string = "impl_sems "
	poolExecs      string = "pool_execs "
	poolStructVars string = "pool_struct_vars "
	poolLinearVars string = "pool_linear_vars "
)

type queryBuilder interface {
	selectRecByRef(implsem.SemRefDS) (string, []any)
}
