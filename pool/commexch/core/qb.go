package core

type QueryBuilder interface {
	InsertRec(ExchRecDS) (string, []any)
	UpdateRec(ExchModDS) (string, []any)
	SelectSnap(ExchQryDS) (string, []any)
}
