package poolconn

const (
	commSems  = "comm_sems"
	poolConns = "pool_conns"
	poolSteps = "pool_steps"
)

type queryBuilder interface {
	insertRec(connRecDS) (string, []any)
	updateRec(connModDS) (string, []any)
	selectSnap(connQryDS) (string, []any)
}
