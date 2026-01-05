package poolexec

import "orglang/orglang/adt/procexp"

type ExecSpecME struct {
	SigQN   string   `json:"sig_qn"`
	ProcIDs []string `json:"proc_ids"`
	SupID   string   `json:"sup_id"`
}

type IdentME struct {
	PoolID string `json:"id" param:"id"`
}

type ExecRefME struct {
	PoolID string `json:"pool_id"`
	ProcID string `json:"proc_id"`
}

type ExecSnapME struct {
	PoolID string      `json:"pool_id"`
	Title  string      `json:"title"`
	Subs   []ExecRefME `json:"subs"`
}

type StepSpecME struct {
	PoolID string            `json:"pool_id"`
	ProcID string            `json:"proc_id"`
	Term   procexp.ExpSpecME `json:"term"`
}
