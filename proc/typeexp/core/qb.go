package core

type QueryBuilder interface {
	InsertRec(StateDS) (string, []any)
	SelectRecByVK() string
}
