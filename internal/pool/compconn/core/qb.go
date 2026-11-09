package core

import connact "orglang/go-engine/internal/pool/connact/core"

type QueryBuilder interface {
	InsertRec(ConnRecDS) (string, []any)
	InsertActRec(connact.ActRecDS) (string, []any)
	UpdateRec(ConnModDS) (string, []any)
	SelectSnap(ConnQryDS) (string, []any)
}
