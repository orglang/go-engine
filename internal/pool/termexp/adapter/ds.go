package adapter

import (
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

type CompSpec struct {
	K       ExpKind     `json:"k"`
	Acquire *GrantSpec  `json:"acquire,omitempty"`
	Accept  *GrantSpec  `json:"accept,omitempty"`
	Hire    *CoopSpec   `json:"hire,omitempty"`
	Apply   *CoopSpec   `json:"apply,omitempty"`
	Release *RevokeSpec `json:"release,omitempty"`
	Detach  *RevokeSpec `json:"detach,omitempty"`
}

type GrantSpec struct {
	ConnChnlPH string   `json:"ph"`
	ContExp    CompSpec `json:"exp"`
}

type CoopSpec struct {
	ConnChnlPH string   `json:"ph"`
	ProcTermQN string   `json:"qn"`
	ContExp    CompSpec `json:"exp"`
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
	ContChnlPH string   `json:"ph"`
	ContExp    CompSpec `json:"exp"`
}

type CoopRec struct {
	ContChnlPH string   `json:"ph"`
	ProcTermQN string   `json:"qn"`
	ContExp    CompSpec `json:"exp"`
}

type RevokeRec struct {
	ContChnlPH string `json:"ph"`
}

// DescData is the adapter representation of a complete stored Definition.
type DescData struct {
	ExpCK  int64       `db:"exp_ck"`
	States []DescState `db:"states"`
}

// DescState is one persisted expression node. Spec contains scalar data and
// child content keys; tree reconstruction belongs to this adapter.
type DescState struct {
	ExpCK    int64     `db:"exp_ck"`
	SupExpCK int64     `db:"sup_exp_ck"`
	K        ExpKind   `db:"kind"`
	Spec     StateSpec `db:"spec"`
}

type StateSpec struct {
	ConnChnlPH  string   `json:"conn_chnl_ph,omitempty"`
	ProcTermQN  string   `json:"proc_term_qn,omitempty"`
	ProcCompQNs []string `json:"proc_comp_qns,omitempty"`
	ContExpCK   int64    `json:"cont_exp_ck,omitempty"`
}
