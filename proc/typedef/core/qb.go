package core

type QueryBuilder interface {
	InsertRec(DefRecDS) (string, []any)
	SelectRecByQN() string
}
