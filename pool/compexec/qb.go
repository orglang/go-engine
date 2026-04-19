package compexec

import (
	"orglang/go-engine/adt/compsem"
)

const (
	implBinds      string = "pool_impl_binds bind"
	compExecs      string = "pool_comp_execs "
	poolStructVars string = "pool_struct_vars "
	poolLinearVars string = "pool_linear_vars "
)

type queryBuilder interface {
	insertRec(execRecDS) (string, []any)
	selectRecByRef(compsem.SemRefDS) (string, []any)
	selectSnapByQN(string) (string, []any)
}
