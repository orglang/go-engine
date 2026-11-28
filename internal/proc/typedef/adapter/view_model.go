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
	ProgID identity.ADT `json:"ProgID"`
	ProgRN seqnum.ADT   `json:"ProgRN"`
}

type DefSnapVP struct {
	TypeRef TypeRefVP `json:"ref"`
	DefSpec DefSpecVP `json:"spec"`
}

func ViewFromDefSnap(snap typedef.DefSnap) DefSnapVP {
	var view DefSnapVP
	view.TypeRef.TypeID = snap.TypeID
	view.TypeRef.ProgID = snap.ProgRef.ProgID
	view.TypeRef.ProgRN = snap.ProgRef.ProgRN
	view.DefSpec.TypeQN = uniqsym.ConvertToString(snap.DefSpec.TypeQN)
	return view
}
