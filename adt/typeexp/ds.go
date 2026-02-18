package typeexp

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descexec"
	"orglang/go-engine/adt/valkey"
)

type Repo interface {
	InsertRec(db.Source, ExpRec, descexec.ExecRef) error
	SelectRecByVK(db.Source, valkey.ADT) (ExpRec, error)
	SelectRecsByVKs(db.Source, []valkey.ADT) ([]ExpRec, error)
	SelectEnv(db.Source, []valkey.ADT) (map[valkey.ADT]ExpRec, error)
}

type expKindDS int

const (
	nonExp expKindDS = iota
	oneExp
	linkExp
	tensorExp
	lolliExp
	plusExp
	withExp
)

type expRefDS struct {
	ExpVK int64     `db:"exp_vk" json:"exp_vk"`
	K     expKindDS `db:"kind" json:"kind"`
}

type expRecDS struct {
	ExpVK  int64
	States []stateDS
}

type stateDS struct {
	ExpVK    int64     `db:"exp_vk"`
	SupExpVK int64     `db:"sup_exp_vk"`
	K        expKindDS `db:"kind"`
	Spec     expSpecDS `db:"spec"`
}

type expSpecDS struct {
	Link   string  `json:"link,omitempty"`
	Tensor *prodDS `json:"tensor,omitempty"`
	Lolli  *prodDS `json:"lolli,omitempty"`
	Plus   []sumDS `json:"plus,omitempty"`
	With   []sumDS `json:"with,omitempty"`
}

type prodDS struct {
	ValExpVK  int64 `json:"on"`
	ContExpVK int64 `json:"to"`
}

type sumDS struct {
	LabQN     string `json:"on"`
	ContExpVK int64  `json:"to"`
}
