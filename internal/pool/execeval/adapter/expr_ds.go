package adapter

import (
	compexp "orglang/go-engine/internal/pool/execeval/core"
)

type ExprSpec struct {
	K       ExprKind     `json:"k"`
	Acquire *GrantSpec  `json:"acquire,omitempty"`
	Accept  *GrantSpec  `json:"accept,omitempty"`
	Hire    *CoopSpec   `json:"hire,omitempty"`
	Apply   *CoopSpec   `json:"apply,omitempty"`
	Release *RevokeSpec `json:"release,omitempty"`
	Detach  *RevokeSpec `json:"detach,omitempty"`
}

type GrantSpec struct {
	CommChnlPH string  `json:"ph"`
	ContExp    ExprSpec `json:"expr"`
}

type CoopSpec struct {
	CommChnlPH string  `json:"ph"`
	ProcTermQN string  `json:"qn"`
	ContExp    ExprSpec `json:"expr"`
}

type RevokeSpec struct {
	CommChnlPH string `json:"ph"`
}

type ExprRec struct {
	K       ExprKind    `json:"k"`
	Acquire *GrantRec  `json:"acquire,omitempty"`
	Accept  *GrantRec  `json:"accept,omitempty"`
	Hire    *CoopRec   `json:"hire,omitempty"`
	Apply   *CoopRec   `json:"apply,omitempty"`
	Release *RevokeRec `json:"release,omitempty"`
	Detach  *RevokeRec `json:"detach,omitempty"`
}

type ExprKind = compexp.ExprKind

const (
	UnkKind     = compexp.UnkKind
	AcquireKind = compexp.AcquireKind
	AcceptKind  = compexp.AcceptKind
	HireKind    = compexp.HireKind
	ApplyKind   = compexp.ApplyKind
	ReleaseKind = compexp.ReleaseKind
	DetachKind  = compexp.DetachKind
)

type GrantRec struct {
	ContChnlPH string  `json:"ph"`
	ContExp    ExprSpec `json:"expr"`
}

type CoopRec struct {
	ContChnlPH string  `json:"ph"`
	ProcTermQN string  `json:"qn"`
	ContExp    ExprSpec `json:"expr"`
}

type RevokeRec struct {
	ContChnlPH string `json:"ph"`
}
