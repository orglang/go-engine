package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/uniqsym"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type DefSpecVP struct {
	TypeQN string `form:"type_qn" json:"type_qn"`
}

type TypeRefVP struct {
	TypeID identity.ADT `json:"TypeID"`
	DescID identity.ADT `json:"DescID"`
	DescRN seqnum.ADT   `json:"DescRN"`
}

type DefSnapVP struct {
	TypeRef TypeRefVP `json:"ref"`
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
