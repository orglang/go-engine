package typedef

import (
	"github.com/orglang/go-sdk/adt/descsem"
)

type DefSpecVP struct {
	TypeQN string `form:"type_qn" json:"type_qn"`
}

type DefSnapVP struct {
	DescRef descsem.SemRef `json:"ref"`
	DefSpec DefSpecVP      `json:"spec"`
}
