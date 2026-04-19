package typedef

const (
	descBinds string = "pool_desc_binds "
	typeDefs  string = "pool_type_defs "
)

type queryBuilder interface {
	insertRec(defRecDS) (string, []any)
	selectRecByQN() string
}
