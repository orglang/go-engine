package adapter

type QueryBuilder interface {
	InsertRec(DefRecDS) (string, []any)
	SelectRecByQN() string
}
