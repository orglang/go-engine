package adapter

import enttypeschema "orglang/go-engine/internal/proc/ent/schema"

type ExpKind int

const (
	UnkKind ExpKind = iota
	OneKind
	LinkKind
	TensorKind
	LolliKind
	PlusKind
	WithKind
)

type ExpRef struct {
	ExpCK int64   `db:"expr_ck" json:"exp_ck"`
	K     ExpKind `db:"kind" json:"kind"`
}

type ExpRec struct {
	ExpCK  int64
	States []State
}

type State struct {
	ExpCK    int64   `db:"expr_ck"`
	SupExprCK int64   `db:"sup_expr_ck"`
	K        ExpKind `db:"kind"`
	Spec     ExpSpec `db:"spec"`
}

type (
	ExpSpec = enttypeschema.ExpSpec
	Prod    = enttypeschema.ExpProd
	Sum     = enttypeschema.ExpSum
)
