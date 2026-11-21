package adapter

import compsemadapter "orglang/go-engine/internal/adt/compsem/adapter"

const (
	implBinds      string = "pool_impl_binds "
	compExecs      string = "pool_comp_execs "
	poolStructVars string = "pool_struct_vars "
	poolLinearVars string = "pool_linear_vars "
)

type queryBuilder interface {
	InsertRec(ExecRec) (string, []any)
	SelectRecByRef(compsemadapter.SemRef) (string, []any)
	SelectSnapByQN(string) (string, []any)
}
