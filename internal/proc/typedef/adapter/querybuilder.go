package adapter

type QueryBuilder interface {
	InsertRec(DefRec) (string, []any)
	SelectRecByQN() string
}
