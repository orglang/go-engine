package core

type QueryBuilder interface {
	InsertRec(TurnRecDS) (string, []any)
}
