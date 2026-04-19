package typedef

const (
	descBinds string = "proc_desc_binds "
	typeDefs  string = "proc_type_defs "
)

type queryBuilder interface {
	insertRec(defRecDS) (string, []any)
	selectRecByQN() string
}
