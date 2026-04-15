package compexec

import (
	"orglang/go-engine/adt/semcomp"
)

const (
	implSems       string = "impl_sems "
	poolExecs      string = "pool_execs "
	poolStructVars string = "pool_struct_vars "
	poolLinearVars string = "pool_linear_vars "
)

type queryBuilder interface {
	selectRecByRef(semcomp.CompRefDS) (string, []any)
}
