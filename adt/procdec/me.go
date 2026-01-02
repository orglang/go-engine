package procdec

import (
	"orglang/orglang/adt/termctx"
)

type SigSpecME struct {
	X     termctx.BindClaimME   `json:"x"`
	SigQN string                `json:"sig_qn"`
	Ys    []termctx.BindClaimME `json:"ys"`
}

type IdentME struct {
	SigID string `json:"id" param:"id"`
}

type SigRefME struct {
	SigID string `json:"id" param:"id"`
	Title string `json:"title"`
	SigRN int64  `json:"rev"`
}

type SigSnapME struct {
	X     termctx.BindClaimME   `json:"x"`
	SigID string                `json:"sig_id"`
	Ys    []termctx.BindClaimME `json:"ys"`
	Title string                `json:"title"`
	SigRN int64                 `json:"sig_rn"`
}
