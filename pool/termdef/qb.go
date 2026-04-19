package termdef

const (
	sigBinds string = "pool_desc_binds "
	poolDecs string = "pool_term_decs "
)

type queryBuilder interface {
	insertRec(decRecDS) (string, []any)
	selectRecByQN(string) (string, []any)
}
