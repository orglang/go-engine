package compexec

const (
	implBinds      string = "proc_impl_binds "
	compExecs      string = "proc_comp_execs "
	procStructVars string = "proc_struct_vars "
	procLinearVars string = "proc_linear_vars "
)

type queryBuilder interface {
	insertRec(execRecDS) (string, []any)
}
