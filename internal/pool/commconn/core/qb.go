package core

import commturn "orglang/go-engine/internal/pool/commturn/core"

type QueryBuilder interface {
	InsertRec(ConnRecDS) (string, []any)
	InsertTurnRec(commturn.TurnRecDS) (string, []any)
	UpdateRec(ConnModDS) (string, []any)
	SelectSnap(ConnQryDS) (string, []any)
}
