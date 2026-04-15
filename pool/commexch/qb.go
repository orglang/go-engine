package commexch

const (
	commSems  = "comm_sems "
	poolConns = "pool_conns "
	poolSteps = "pool_steps "
)

type queryBuilder interface {
	insertRec(exchRecDS) (string, []any)
	updateRec(exchModDS) (string, []any)
	selectSnap(exchQryDS) (string, []any)
}
