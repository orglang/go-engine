package adapter

type queryBuilder interface {
	insertRec(semRecDS) (string, []any)
}
