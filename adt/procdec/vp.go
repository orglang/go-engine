package procdec

import (
	"github.com/orglang/go-sdk/adt/descsem"
)

type DecSpecVP struct {
	ProcQN string `form:"qn" json:"qn"`
}

type DecSnapVP struct {
	DescRef descsem.SemRef `json:"ref"`
}
