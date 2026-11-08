package core

type QueryBuilder interface {
	InsertRec(ConnRecDS) (string, []any)
	UpdateRec(ConnModDS) (string, []any)
	SelectSnap(ConnQryDS) (string, []any)
}
