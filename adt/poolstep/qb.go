package poolstep

const (
	commSems  = "comm_sems"
	poolSteps = "pool_steps"
)

type queryBuilder interface {
	insertRec(StepRecDS) (string, []any)
}
