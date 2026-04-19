package termdec

import (
	"github.com/orglang/go-sdk/adt/termsem"
)

type DecSpecVP struct {
	TermQN string `form:"qn" json:"qn"`
}

type DecSnapVP struct {
	TermRef termsem.SemRef `json:"ref"`
}
