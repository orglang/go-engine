package typedef

import (
	"github.com/orglang/go-sdk/adt/uniqref"
)

type DefSpecVP struct {
	TypeQN string `form:"type_qn" json:"type_qn"`
}

type DefRefVP = uniqref.Msg

type DefSnapVP struct {
	DefRef  DefRefVP  `json:"ref"`
	DefSpec DefSpecVP `json:"spec"`
}
