package commexch

const (
	commExchs = "pool_comm_exchs "
	commTurns = "pool_comm_turns "
)

type queryBuilder interface {
	insertRec(exchRecDS) (string, []any)
	updateRec(exchModDS) (string, []any)
	selectSnap(exchQryDS) (string, []any)
}
