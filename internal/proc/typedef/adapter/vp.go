package adapter

import (
	"orglang/go-engine/internal/adt/typesem"
)

type DefSpecVP struct {
	TypeQN string `form:"type_qn" json:"type_qn"`
}

type DefSnapVP struct {
	TypeRef typesem.SemRef `json:"ref"`
	DefSpec DefSpecVP      `json:"spec"`
}
