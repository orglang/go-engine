package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

type DecSpecVP struct { TermQN string `form:"qn" json:"qn"` }

type DecSnapVP struct {
	TermRef struct {
		TermID identity.ADT `json:"TermID"`
		DescID identity.ADT `json:"DescID"`
		DescRN seqnum.ADT   `json:"DescRN"`
	} `json:"ref"`
}
