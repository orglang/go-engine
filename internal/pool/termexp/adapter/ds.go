package adapter

import termexp "orglang/go-engine/internal/pool/termexp/core"

type ExpSpec struct {
	K       ExpKind     `json:"k"`
	Acquire *GrantSpec  `json:"acquire,omitempty"`
	Accept  *GrantSpec  `json:"accept,omitempty"`
	Hire    *CoopSpec   `json:"hire,omitempty"`
	Apply   *CoopSpec   `json:"apply,omitempty"`
	Release *RevokeSpec `json:"release,omitempty"`
	Detach  *RevokeSpec `json:"detach,omitempty"`
}

type GrantSpec struct {
	ConnChnlPH string  `json:"ph"`
	ContExp    ExpSpec `json:"exp"`
}

type CoopSpec struct {
	ConnChnlPH string  `json:"ph"`
	ProcTermQN string  `json:"qn"`
	ContExp    ExpSpec `json:"exp"`
}

type RevokeSpec struct {
	ConnChnlPH string `json:"ph"`
}

type ExpRec struct {
	K       ExpKind    `json:"k"`
	Acquire *GrantRec  `json:"acquire,omitempty"`
	Accept  *GrantRec  `json:"accept,omitempty"`
	Hire    *CoopRec   `json:"hire,omitempty"`
	Apply   *CoopRec   `json:"apply,omitempty"`
	Release *RevokeRec `json:"release,omitempty"`
	Detach  *RevokeRec `json:"detach,omitempty"`
}

type ExpKind = termexp.ExpKind

const (
	UnkKind     = termexp.UnkKind
	AcquireKind = termexp.AcquireKind
	AcceptKind  = termexp.AcceptKind
	HireKind    = termexp.HireKind
	ApplyKind   = termexp.ApplyKind
	ReleaseKind = termexp.ReleaseKind
	DetachKind  = termexp.DetachKind
)

type GrantRec struct {
	ContChnlPH string  `json:"ph"`
	ContExp    ExpSpec `json:"exp"`
}

type CoopRec struct {
	ContChnlPH string  `json:"ph"`
	ProcTermQN string  `json:"qn"`
	ContExp    ExpSpec `json:"exp"`
}

type RevokeRec struct {
	ContChnlPH string `json:"ph"`
}
