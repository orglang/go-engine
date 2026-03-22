package poolexec

import (
	"orglang/go-engine/adt/implsem"
)

const (
	implSems       string = "impl_sems sem"
	implBinds      string = "impl_binds bind"
	poolExecs      string = "pool_execs exec"
	poolVars       string = "pool_vars var"
	poolStructVars string = "pool_struct_vars struct_var"
	poolLinearVars string = "pool_linear_vars linear_var"
)

type queryBuilder interface {
	insertRec(execRecDS) (string, []any)
	selectSnap(implsem.SemRefDS) (string, []any)
	selectSnapByQN(string) (string, []any)
}
