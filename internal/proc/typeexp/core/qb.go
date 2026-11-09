package core

type QueryBuilder interface {
	InsertRec(StateDS) (string, []any)
	SelectRecByCK() string
}
