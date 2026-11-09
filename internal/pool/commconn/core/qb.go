package core

import commact "orglang/go-engine/internal/pool/commact/core"

type QueryBuilder interface {
	InsertRec(ConnRecDS) (string, []any)
	InsertActRec(commact.ActRecDS) (string, []any)
	UpdateRec(ConnModDS) (string, []any)
	SelectSnap(ConnQryDS) (string, []any)
}
