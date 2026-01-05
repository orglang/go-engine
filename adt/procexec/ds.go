package procexec

import (
	"database/sql"

	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/procexp"
)

type execRepo interface {
	SelectMain(sd.Source, identity.ADT) (MainCfg, error)
	UpdateMain(sd.Source, MainMod) error
}

type SemRepo interface {
	InsertSem(sd.Source, ...SemRec) error
	SelectSemByID(sd.Source, identity.ADT) (SemRec, error)
}

type modDS struct {
	Locks []lockDS
	Bnds  []bndDS
	Steps []SemRecDS
}

type lockDS struct {
	PoolID string
	PoolRN int64
}

type bndDS struct {
	ProcID  string
	ChnlPH  string
	ChnlID  string
	StateID string
	PoolRN  int64
}

type SemRecDS struct {
	ID  string           `db:"id"`
	K   semKindDS        `db:"kind"`
	PID sql.NullString   `db:"pid"`
	VID sql.NullString   `db:"vid"`
	TR  procexp.ExpRecDS `db:"spec"`
}

type semKindDS int

const (
	nonsem = semKindDS(iota)
	msgKind
	svcKind
)
