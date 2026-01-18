package procdec

import "github.com/orglang/go-sdk/adt/uniqref"

type DecRefVP = uniqref.Msg

type DecSpecVP struct {
	ProcNS string `form:"ns" json:"ns"`
	ProcSN string `form:"sn" json:"sn"`
	Title  string `form:"name" json:"title"`
}

type DecSnapVP struct {
	DecID string `json:"dec_id"`
	DecRN int64  `json:"dec_rn"`
	Title string `json:"title"`
}
