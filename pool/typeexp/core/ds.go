package core

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/typesem"
	"orglang/go-engine/adt/valkey"
)

type Repo interface {
	AddRec(db.UoW, ExpRec, typesem.SemRef) error
	GetRecByVK(db.UoW, valkey.ADT) (ExpRec, error)
	GetRecsByVKs(db.UoW, []valkey.ADT) ([]ExpRec, error)
	GetRecMap(db.UoW, map[symbol.ADT]valkey.ADT) (map[symbol.ADT]ExpRec, error)
}

type ExpKind int16

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

type ExpRefDS struct {
	ExpVK int64   `db:"exp_vk" json:"exp_vk"`
	K     ExpKind `db:"kind" json:"kind"`
}

type ExpRecDS struct {
	ExpVK  int64
	States []StateDS
}

type StateDS struct {
	ExpVK    int64     `db:"exp_vk"`
	SupExpVK int64     `db:"sup_exp_vk"`
	K        ExpKind   `db:"kind"`
	Spec     ExpSpecDS `db:"spec" fieldopt:"noexpand"`
}

type ExpSpecDS struct {
	Link   string   `json:"link,omitempty"`
	Tensor *ProdDS  `json:"tensor,omitempty"`
	Lolli  *ProdDS  `json:"lolli,omitempty"`
	Plus   *SumDS   `json:"plus,omitempty"`
	With   *SumDS   `json:"with,omitempty"`
	Up     *ShiftDS `json:"up,omitempty"`
	Down   *ShiftDS `json:"down,omitempty"`
}

type ProdDS struct {
	ValExpVK  int64 `json:"on"`
	ContExpVK int64 `json:"to"`
}

type SumDS struct {
	ProcQNs   []string `json:"on"`
	ContExpVK int64    `json:"to"`
}

type ShiftDS struct {
	ContExpVK int64 `json:"to"`
}
