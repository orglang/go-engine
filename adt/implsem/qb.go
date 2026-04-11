package implsem

const (
	implSems string = "impl_sems"
)

type queryBuilder interface {
	insertRec(semRecDS) (string, []any)
	updateRec(SemRefDS) (string, []any)
}
