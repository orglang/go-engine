package core

const (
	implBinds      string = "proc_impl_binds "
	CompExecs      string = "proc_comp_execs "
	procStructVars string = "proc_struct_vars "
	procLinearVars string = "proc_linear_vars "
)

type QueryBuilder interface {
	InsertRec(ExecRecDS) (string, []any)
}
