package adapter

import (
	termdec "orglang/go-engine/internal/proc/termdec/core"
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

func ViewFromDecSnap(snap termdec.DecSnap) DecSnapVP {
	var view DecSnapVP
	view.TermRef.TermID = snap.TermID
	view.TermRef.DescID = snap.DescRef.DescID
	view.TermRef.DescRN = snap.DescRef.DescRN
	return view
}
