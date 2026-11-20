package adapter

import compexec "orglang/go-engine/internal/proc/compexec/core"

const (
	implBinds      string = "proc_impl_binds "
	CompExecs      string = "proc_comp_execs "
	procStructVars string = "proc_struct_vars "
	procLinearVars string = "proc_linear_vars "
)

type queryBuilder interface {
	InsertRec(compexec.ExecRecDS) (string, []any)
}
