package poolstep

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/poolexp"
)

type Repo interface {
	AddRec(db.Source, StepRec) error
	AddRecs(db.Source, []StepRec) error
}

type StepRecDS struct {
	CommID string           `db:"comm_id"`
	CommRN int64            `db:"comm_rn"`
	ImplID string           `db:"impl_id"`
	ChnlID string           `db:"chnl_id"`
	K      stepKind         `db:"kind"`
	Exp    poolexp.ExpRecDS `db:"exp" fieldopt:"noexpand"`
}

type stepKind int16

const (
	unkKind stepKind = iota
	PubKind
	SubKind
)
