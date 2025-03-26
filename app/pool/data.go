package pool

import (
	"database/sql"

	"smecalculus/rolevod/internal/proc"
)

type refData struct {
	PoolID string `json:"pool_id"`
	Title  string `json:"title"`
}

type subSnapData struct {
	PoolID string    `db:"pool_id"`
	Title  string    `db:"title"`
	Subs   []refData `db:"subs"`
}

type rootData struct {
	PoolID string         `db:"pool_id"`
	ProcID string         `db:"proc_id"`
	Title  string         `db:"title"`
	SupID  sql.NullString `db:"sup_pool_id"`
	Rev    int64          `db:"rev"`
}

type liabData struct {
	PoolID string `db:"pool_id"`
	ProcID string `db:"proc_id"`
	Rev    int64  `db:"rev"`
}

type epData struct {
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

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	DataToRef       func(refData) (Ref, error)
	DataFromRef     func(Ref) refData
	DataToRefs      func([]refData) ([]Ref, error)
	DataFromRefs    func([]Ref) []refData
	DataToRoot      func(rootData) (Root, error)
	DataFromRoot    func(Root) rootData
	DataToLiab      func(liabData) (proc.Liab, error)
	DataFromLiab    func(proc.Liab) liabData
	DataToSubSnap   func(subSnapData) (SubSnap, error)
	DataFromSubSnap func(SubSnap) subSnapData
	DataToEPs       func([]epData) ([]proc.EP, error)
)
