package adapter

type queryBuilder interface {
	insertRec(semRec) (string, []any)
}
