package procexec

import (
	"orglang/orglang/adt/procdef"
)

type SpecME struct {
	ProcID string             `json:"proc_id" param:"id"`
	PoolID string             `json:"pool_id"`
	Term   procdef.CallSpecME `json:"term"`
}

type IdentME struct {
	ProcID string `param:"id"`
}

type RefME struct {
	ProcID string `json:"proc_id"`
}

type SnapME struct {
	ProcID string `json:"proc_id"`
}
