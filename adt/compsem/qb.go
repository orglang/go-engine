package compsem

type queryBuilder interface {
	updateRef(SemRefDS) (string, []any)
}
