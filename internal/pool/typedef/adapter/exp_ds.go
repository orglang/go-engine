package adapter

import enttypeschema "orglang/go-engine/internal/pool/ent/schema"

type ExpKind int

const (
	UnkKind ExpKind = iota
	OneKind
	LinkKind
	TensorKind
	LolliKind
	PlusKind
	WithKind
	UpKind
	DownKind
)

type ExpRef struct {
	ExpCK int64   `db:"exp_ck" json:"exp_ck"`
	K     ExpKind `db:"kind" json:"kind"`
}

type ExpRec struct {
	ExpCK  int64
	States []State
}

type State struct {
	ExpCK    int64   `db:"exp_ck"`
	SupExpCK int64   `db:"sup_exp_ck"`
	K        ExpKind `db:"kind"`
	Spec     ExpSpec `db:"spec" fieldopt:"noexpand"`
}

// Keep the adapter DTO aligned with the JSON type used by the Ent schema.
type (
	ExpSpec = enttypeschema.ExpSpec
	Prod    = enttypeschema.ExpProd
	Sum     = enttypeschema.ExpSum
	Shift   = enttypeschema.ExpShift
)
