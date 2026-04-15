package termdec

const (
	descBinds string = "desc_binds"
	poolDecs  string = "pool_decs"
)

type queryBuilder interface {
	insertRec(decRecDS) (string, []any)
	selectRecByQN(string) (string, []any)
}
