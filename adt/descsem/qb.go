package descsem

type queryBuilder interface {
	insertRec(SemRecDS) (string, []any)
}
