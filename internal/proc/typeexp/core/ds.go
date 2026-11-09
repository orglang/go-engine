package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ExpRec) error
	SelectRecByCK(context.Context, tx.Conn, contkey.ADT) (ExpRec, error)
	SelectRecsByCKs(context.Context, tx.Conn, []contkey.ADT) ([]ExpRec, error)
	SelectEnv(context.Context, tx.Conn, []contkey.ADT) (map[contkey.ADT]ExpRec, error)
}

type ExpKindDS int

const (
	UnkKind ExpKindDS = iota
	OneKind
	LinkKind
	TensorKind
	LolliKind
	PlusKind
	WithKind
)

type ExpRefDS struct {
	ExpCK int64     `db:"exp_ck" json:"exp_ck"`
	K     ExpKindDS `db:"kind" json:"kind"`
}

type ExpRecDS struct {
	ExpCK  int64
	States []StateDS
}

type StateDS struct {
	ExpCK    int64     `db:"exp_ck"`
	SupExpCK int64     `db:"sup_exp_ck"`
	K        ExpKindDS `db:"kind"`
	Spec     ExpSpecDS `db:"spec"`
}

type ExpSpecDS struct {
	Link   string  `json:"link,omitempty"`
	Tensor *ProdDS `json:"tensor,omitempty"`
	Lolli  *ProdDS `json:"lolli,omitempty"`
	Plus   []SumDS `json:"plus,omitempty"`
	With   []SumDS `json:"with,omitempty"`
}

type ProdDS struct {
	ValExpCK  int64 `json:"on"`
	ContExpCK int64 `json:"to"`
}

type SumDS struct {
	LabQN     string `json:"on"`
	ContExpCK int64  `json:"to"`
}
