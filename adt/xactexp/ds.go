package xactexp

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
)

type Repo interface {
	InsertRec(db.Source, ExpRec) error
	SelectRecByID(db.Source, identity.ADT) (ExpRec, error)
	SelectRecsByIDs(db.Source, []identity.ADT) ([]ExpRec, error)
	SelectEnv(db.Source, []identity.ADT) (map[identity.ADT]ExpRec, error)
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
