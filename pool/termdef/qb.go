package termdef

const (
	descBinds string = "pool_desc_binds "
	termDefs  string = "pool_term_decs "
)

type queryBuilder interface {
	insertRec(defRecDS) (string, []any)
	selectRecByQN(string) (string, []any)
}
