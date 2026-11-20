package adapter

import (
	"orglang/go-engine/internal/adt/compsem"
	compexec "orglang/go-engine/internal/pool/compexec/core"
)

const (
	implBinds      string = "pool_impl_binds "
	compExecs      string = "pool_comp_execs "
	poolStructVars string = "pool_struct_vars "
	poolLinearVars string = "pool_linear_vars "
)

type queryBuilder interface {
	InsertRec(compexec.ExecRecDS) (string, []any)
	SelectRecByRef(compsem.SemRefDS) (string, []any)
	SelectSnapByQN(string) (string, []any)
}
