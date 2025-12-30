package poolexec

import (
	"orglang/orglang/adt/procdef"
)

type PoolSpecME struct {
	SigQN   string   `json:"sig_qn"`
	ProcIDs []string `json:"proc_ids"`
	SupID   string   `json:"sup_id"`
}

type IdentME struct {
	PoolID string `json:"id" param:"id"`
}

type PoolRefME struct {
	PoolID string `json:"pool_id"`
	ProcID string `json:"proc_id"`
}

type PoolSnapME struct {
	PoolID string      `json:"id"`
	Title  string      `json:"title"`
	Subs   []PoolRefME `json:"subs"`
}

type StepSpecME struct {
	PoolID string             `json:"pool_id"`
	ProcID string             `json:"proc_id"`
	Term   procdef.TermSpecME `json:"term"`
}
