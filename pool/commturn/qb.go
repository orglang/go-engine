package commturn

const (
	commSems  = "comm_sems"
	poolSteps = "pool_steps"
)

type queryBuilder interface {
	insertRec(TurnRecDS) (string, []any)
}
