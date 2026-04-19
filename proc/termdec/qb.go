package termdec

const (
	descBinds string = "proc_desc_binds "
	termDecs  string = "proc_term_decs "
	typeDefs  string = "proc_type_defs "
)

type queryBuilder interface {
	insertRec(decRecDS) (string, []any)
}
