package typedef

import (
	"github.com/orglang/go-sdk/adt/typeexp"
	"github.com/orglang/go-sdk/adt/uniqref"
)

type DefSpecVP struct {
	TypeNS string `form:"ns" json:"ns"`
	TypeSN string `form:"name" json:"name"`
}

type DefRefVP = uniqref.Msg

type DefSnapVP struct {
	DefID  string            `json:"def_id"`
	DefRN  int64             `json:"def_rn"`
	Title  string            `json:"title"`
	TypeES typeexp.ExpSpecME `json:"type_es"`
}
