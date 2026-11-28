package adapter

import enttypeschema "orglang/go-engine/internal/proc/ent/schema"

type ExprKind int

const (
	UnkKind ExprKind = iota
	OneKind
	LinkKind
	TensorKind
	LolliKind
	PlusKind
	WithKind
)

type ExprRef struct {
	ExprCK int64   `db:"expr_ck" json:"expr_ck"`
	K     ExprKind `db:"kind" json:"kind"`
}

type ExprRec struct {
	ExprCK  int64
	States []State
}

type State struct {
	ExprCK     int64   `db:"expr_ck"`
	SupExprCK int64   `db:"sup_expr_ck"`
	K         ExprKind `db:"kind"`
	Spec      ExprSpec `db:"spec"`
}

type (
	ExprSpec = enttypeschema.ExprSpec
	Prod    = enttypeschema.ExprProd
	Sum     = enttypeschema.ExprSum
)
