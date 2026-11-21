package adapter

import connactadapter "orglang/go-engine/internal/pool/connact/adapter"

type QueryBuilder interface {
	InsertRec(ConnRec) (string, []any)
	InsertActRec(connactadapter.ActRec) (string, []any)
	UpdateRec(ConnMod) (string, []any)
	SelectSnap(ConnQry) (string, []any)
}
