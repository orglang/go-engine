package poolexec

import (
	"orglang/go-engine/adt/implsem"
)

const (
	execs      = "pool_execs exec"
	structVars = "pool_struct_vars var"
	linearVars = "pool_linear_vars var"
	sems       = "impl_sems sem"
)

type queryBuilder interface {
	insertRec(execRecDS) (string, []any)
	selectSnap(implsem.SemRefDS) (string, []any)
}
