package commturn

const (
	commTurns = "pool_comm_turns "
)

type queryBuilder interface {
	insertRec(TurnRecDS) (string, []any)
}
