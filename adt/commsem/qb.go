package commsem

const (
	commSems string = "comm_sems sem"
)

type queryBuilder interface {
	insertRec(semRecDS) (string, []any)
}
