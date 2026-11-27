package adapter

import (
	typedef "orglang/go-engine/internal/proc/typedef/core"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/uniqsym"
)

type DefSpecVP struct {
	TypeQN string `form:"type_qn" json:"type_qn"`
}

type DefSnapVP struct {
	TypeRef struct {
		TypeID identity.ADT `json:"TypeID"`
		DescID identity.ADT `json:"DescID"`
		DescRN seqnum.ADT   `json:"DescRN"`
	} `json:"ref"`
	DefSpec DefSpecVP `json:"spec"`
}

func ViewFromDefSnap(snap typedef.DefSnap) DefSnapVP {
	var view DefSnapVP
	view.TypeRef.TypeID = snap.TypeID
	view.TypeRef.DescID = snap.DescRef.DescID
	view.TypeRef.DescRN = snap.DescRef.DescRN
	view.DefSpec.TypeQN = uniqsym.ConvertToString(snap.DefSpec.TypeQN)
	return view
}
