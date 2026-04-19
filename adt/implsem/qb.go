package implsem

type queryBuilder interface {
	insertRec(SemRecDS) (string, []any)
}
