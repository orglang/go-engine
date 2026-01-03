package typeexp

import (
	"database/sql"

	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
)

type Repo interface {
	Insert(sd.Source, ExpRec) error
	SelectRecByID(sd.Source, identity.ADT) (ExpRec, error)
	SelectRecsByIDs(sd.Source, []identity.ADT) ([]ExpRec, error)
	SelectEnv(sd.Source, []identity.ADT) (map[identity.ADT]ExpRec, error)
}

type expKindDS int

const (
	nonterm expKindDS = iota
	oneKind
	linkKind
	tensorKind
	lolliKind
	plusKind
	withKind
)

type ExpRefDS struct {
	ExpID string    `db:"exp_id" json:"exp_id"`
	K     expKindDS `db:"kind" json:"kind"`
}

type expRecDS struct {
	ExpID  string
	States []stateDS
}

type stateDS struct {
	ExpID  string         `db:"exp_id"`
	K      expKindDS      `db:"kind"`
	FromID sql.NullString `db:"from_id"`
	Spec   expSpecDS      `db:"spec"`
}

type expSpecDS struct {
	Link   string  `json:"link,omitempty"`
	Tensor *prodDS `json:"tensor,omitempty"`
	Lolli  *prodDS `json:"lolli,omitempty"`
	Plus   []sumDS `json:"plus,omitempty"`
	With   []sumDS `json:"with,omitempty"`
}

type prodDS struct {
	ValES  string `json:"on"`
	ContES string `json:"to"`
}

type sumDS struct {
	Lab    string `json:"on"`
	ContES string `json:"to"`
}
