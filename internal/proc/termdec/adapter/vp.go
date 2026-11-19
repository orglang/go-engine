package adapter

import (
	"orglang/go-engine/internal/adt/termsem"
)

type DecSpecVP struct {
	TermQN string `form:"qn" json:"qn"`
}

type DecSnapVP struct {
	TermRef termsem.SemRef `json:"ref"`
}
