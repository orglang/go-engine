package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
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
