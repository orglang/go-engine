package typedef

import (
	"github.com/orglang/go-sdk/adt/typesem"
)

type DefSpecVP struct {
	TypeQN string `form:"type_qn" json:"type_qn"`
}

type DefSnapVP struct {
	TypeRef typesem.SemRef `json:"ref"`
	DefSpec DefSpecVP      `json:"spec"`
}
