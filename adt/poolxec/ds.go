package poolxec

import (
	"database/sql"

	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/procxec"
)

// Port
type Repo interface {
	Insert(sd.Source, ExecRec) error
	InsertLiab(sd.Source, procxec.Liab) error
	SelectRefs(sd.Source) ([]ExecRef, error)
	SelectSubs(sd.Source, identity.ADT) (ExecSnap, error)
	SelectProc(sd.Source, identity.ADT) (procxec.Cfg, error)
	UpdateProc(sd.Source, procxec.Mod) error
}

type execRefDS struct {
	PoolID string `json:"pool_id"`
	ProcID string `json:"proc_id"`
}

type execSnapDS struct {
	PoolID string      `db:"pool_id"`
	Title  string      `db:"title"`
	Subs   []execRefDS `db:"subs"`
}

type execRecDS struct {
	PoolID string         `db:"pool_id"`
	ProcID string         `db:"proc_id"`
	SupID  sql.NullString `db:"sup_pool_id"`
	PoolRN int64          `db:"rev"`
}

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
