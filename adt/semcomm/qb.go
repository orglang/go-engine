package semcomm

const (
	commSems string = "comm_sems"
)

type queryBuilder interface {
	insertRec(semRecDS) (string, []any)
	updateRec(SemRefDS) (string, []any)
}
