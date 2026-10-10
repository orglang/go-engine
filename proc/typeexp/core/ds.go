package core

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/valkey"
)

type Repo interface {
	AddRec(db.UoW, ExpRec) error
	SelectRecByVK(db.UoW, valkey.ADT) (ExpRec, error)
	SelectRecsByVKs(db.UoW, []valkey.ADT) ([]ExpRec, error)
	SelectEnv(db.UoW, []valkey.ADT) (map[valkey.ADT]ExpRec, error)
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
	ExpVK int64     `db:"exp_vk" json:"exp_vk"`
	K     ExpKindDS `db:"kind" json:"kind"`
}

type ExpRecDS struct {
	ExpVK  int64
	States []StateDS
}

type StateDS struct {
	ExpVK    int64     `db:"exp_vk"`
	SupExpVK int64     `db:"sup_exp_vk"`
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
	ValExpVK  int64 `json:"on"`
	ContExpVK int64 `json:"to"`
}

type SumDS struct {
	LabQN     string `json:"on"`
	ContExpVK int64  `json:"to"`
}
