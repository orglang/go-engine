package typedef

import (
	"github.com/orglang/go-sdk/adt/descexec"
)

type DefSpecVP struct {
	TypeQN string `form:"type_qn" json:"type_qn"`
}

type DefSnapVP struct {
	DescRef descexec.ExecRef `json:"ref"`
	DefSpec DefSpecVP        `json:"spec"`
}
