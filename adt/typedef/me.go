package typedef

import (
	"orglang/orglang/adt/typeexp"
)

type DefSpecME struct {
	TypeQN string            `json:"type_qn"`
	TypeES typeexp.ExpSpecME `json:"type_es"`
}

type IdentME struct {
	DefID string `json:"def_id" param:"id"`
}

type DefRefME struct {
	DefID string `json:"def_id" param:"id"`
	DefRN int64  `json:"def_rn" query:"rn"`
}

type DefSnapME struct {
	DefID  string            `json:"def_id" param:"id"`
	DefRN  int64             `json:"def_rn" query:"rn"`
	Title  string            `json:"title"`
	TypeQN string            `json:"type_qn"`
	TypeES typeexp.ExpSpecME `json:"type_es"`
}
