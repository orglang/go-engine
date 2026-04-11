package procexec

const (
	implSems       string = "impl_sems"
	implBinds      string = "impl_binds"
	procExecs      string = "proc_execs"
	procStructVars string = "proc_struct_vars"
	procLinearVars string = "proc_linear_vars"
)

type queryBuilder interface {
	insertRec(execRecDS) (string, []any)
}
