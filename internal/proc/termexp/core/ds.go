package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	Insert(context.Context, tx.Conn, ExpRec) error
}

type ExpSpecDS struct {
	K     ExpKind      `json:"k"`
	Close *CloseSpecDS `json:"close,omitempty"`
	Wait  *WaitSpecDS  `json:"wait,omitempty"`
	Send  *SendSpecDS  `json:"send,omitempty"`
	Recv  *RecvSpecDS  `json:"recv,omitempty"`
	Lab   *LabSpecDS   `json:"lab,omitempty"`
	Case  *CaseSpecDS  `json:"case,omitempty"`
	Fwd   *FwdSpecDS   `json:"fwd,omitempty"`
}

type ExpRecDS struct {
	K     ExpKind     `json:"k"`
	Close *CloseRecDS `json:"close,omitempty"`
	Wait  *WaitRecDS  `json:"wait,omitempty"`
	Send  *SendRecDS  `json:"send,omitempty"`
	Recv  *RecvRecDS  `json:"recv,omitempty"`
	Lab   *LabRecDS   `json:"lab,omitempty"`
	Case  *CaseRecDS  `json:"case,omitempty"`
	Fwd   *FwdRecDS   `json:"fwd,omitempty"`
}

type ExpKind int

const (
	NonExp ExpKind = iota
	CloseExp
	WaitExp
	SendExp
	RecvExp
	LabExp
	CaseExp
	LinkExp
	SpawnExp
	FwdExp
)

type CloseSpecDS struct {
	X string `json:"x"`
}

type CloseRecDS struct {
	X string `json:"x"`
}

type WaitSpecDS struct {
	X      string    `json:"x"`
	ContES ExpSpecDS `json:"cont"`
}

type WaitRecDS struct {
	X      string    `json:"x"`
	ContES ExpSpecDS `json:"cont"`
}

type SendSpecDS struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type SendRecDS struct {
	X        string `json:"x"`
	A        string `json:"a"`
	B        string `json:"b"`
	ConnID   string `json:"conn_id"`
	// ValExpCK keeps the legacy JSON property for wire compatibility.
	ValExpCK int64  `json:"val_exp_vk"`
}

type RecvSpecDS struct {
	X      string    `json:"x"`
	Y      string    `json:"y"`
	ContES ExpSpecDS `json:"cont"`
}

type RecvRecDS struct {
	X      string    `json:"x"`
	A      string    `json:"a"`
	Y      string    `json:"y"`
	ContES ExpSpecDS `json:"cont"`
}

type LabSpecDS struct {
	X     string `json:"x"`
	Label string `json:"lab"`
}

type LabRecDS struct {
	X        string `json:"x"`
	A        string `json:"a"`
	Label    string `json:"lab"`
	// ValExpCK keeps the legacy JSON property for wire compatibility.
	ValExpCK int64  `json:"val_exp_vk"`
}

type CaseSpecDS struct {
	X        string         `json:"x"`
	Branches []BranchSpecDS `json:"brs"`
}

type CaseRecDS struct {
	X        string        `json:"x"`
	A        string        `json:"a"`
	Branches []BranchRecDS `json:"brs"`
}

type BranchSpecDS struct {
	Label  string    `json:"lab"`
	ContES ExpSpecDS `json:"cont"`
}

type BranchRecDS struct {
	Label  string    `json:"lab"`
	ContES ExpSpecDS `json:"cont"`
}

type FwdSpecDS struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type FwdRecDS struct {
	X string `json:"x"`
	B string `json:"b"`
}
