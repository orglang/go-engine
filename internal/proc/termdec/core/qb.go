package core

type QueryBuilder interface {
	InsertRec(DecRecDS) (string, []any)
}
