package procexec

import (
	"orglang/go-runtime/lib/db"

	"orglang/go-runtime/adt/identity"
	"orglang/go-runtime/adt/procbind"
	"orglang/go-runtime/adt/procstep"
	"orglang/go-runtime/adt/uniqref"
)

type Repo interface {
	SelectSnap(db.Source, ExecRef) (ExecSnap, error)
	UpdateProc(db.Source, ExecMod) error
	SelectMain(db.Source, identity.ADT) (MainCfg, error)
	UpdateMain(db.Source, MainMod) error
}

type execModDS struct {
	Locks []execRefDS
	Binds []procbind.BindRecDS
	Steps []procstep.StepRecDS
}

type execRefDS = uniqref.Data

type liabDS struct {
	PoolID string `db:"pool_id"`
	ProcID string `db:"proc_id"`
	PoolRN int64  `db:"rev"`
}

type epDS struct {
	ProcID   string  `db:"proc_id"`
	ChnlPH   string  `db:"chnl_ph"`
	ChnlID   string  `db:"chnl_id"`
	StateID  string  `db:"state_id"`
	PoolID   string  `db:"pool_id"`
	SrvID    string  `db:"srv_id"`
	SrvRevs  []int64 `db:"srv_revs"`
	ClntID   string  `db:"clnt_id"`
	ClntRevs []int64 `db:"clnt_revs"`
}
