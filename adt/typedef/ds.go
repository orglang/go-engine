package typedef

import (
	"database/sql"

	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
)

type Repo interface {
	InsertType(sd.Source, DefRec) error
	UpdateType(sd.Source, DefRec) error
	SelectTypeRefs(sd.Source) ([]DefRef, error)
	SelectTypeRecByID(sd.Source, identity.ADT) (DefRec, error)
	SelectTypeRecsByIDs(sd.Source, []identity.ADT) ([]DefRec, error)
	SelectTypeRecByQN(sd.Source, qualsym.ADT) (DefRec, error)
	SelectTypeRecsByQNs(sd.Source, []qualsym.ADT) ([]DefRec, error)
	SelectTypeEnv(sd.Source, []qualsym.ADT) (map[qualsym.ADT]DefRec, error)

	InsertTerm(sd.Source, TermRec) error
	SelectTermRecByID(sd.Source, identity.ADT) (TermRec, error)
	SelectTermRecsByIDs(sd.Source, []identity.ADT) ([]TermRec, error)
	SelectTermEnv(sd.Source, []identity.ADT) (map[identity.ADT]TermRec, error)
}

type defRefDS struct {
	DefID string `db:"def_id"`
	DefRN int64  `db:"def_rn"`
}

type defRecDS struct {
	DefID  string `db:"def_id"`
	Title  string `db:"title"`
	TermID string `db:"term_id"`
	DefRN  int64  `db:"def_rn"`
}

type termKindDS int

const (
	nonterm termKindDS = iota
	oneKind
	linkKind
	tensorKind
	lolliKind
	plusKind
	withKind
)

type TermRefDS struct {
	TermID string     `db:"term_id" json:"term_id"`
	K      termKindDS `db:"kind" json:"kind"`
}

type termRecDS struct {
	TermID string
	States []stateDS
}

type stateDS struct {
	TermID string         `db:"term_id"`
	K      termKindDS     `db:"kind"`
	FromID sql.NullString `db:"from_id"`
	Spec   termSpecDS     `db:"spec"`
}

type termSpecDS struct {
	Link   string  `json:"link,omitempty"`
	Tensor *prodDS `json:"tensor,omitempty"`
	Lolli  *prodDS `json:"lolli,omitempty"`
	Plus   []sumDS `json:"plus,omitempty"`
	With   []sumDS `json:"with,omitempty"`
}

type prodDS struct {
	Val  string `json:"on"`
	Cont string `json:"to"`
}

type sumDS struct {
	Lab  string `json:"on"`
	Cont string `json:"to"`
}
