package procexec

import (
	"orglang/orglang/adt/procdef"
)

type ExecSpecME struct {
	ProcID string             `json:"exec_id" param:"id"`
	PoolID string             `json:"pool_id"`
	Term   procdef.CallSpecME `json:"term"`
}

type IdentME struct {
	ProcID string `param:"id"`
}

type ExecRefME struct {
	ExecID string `json:"exec_id"`
}

type ExecSnapME struct {
	ExecID string `json:"exec_id"`
}
