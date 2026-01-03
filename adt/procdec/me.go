package procdec

import (
	"orglang/orglang/adt/termctx"
)

type DecSpecME struct {
	X      termctx.BindClaimME   `json:"x"`
	ProcQN string                `json:"proc_qn"`
	Ys     []termctx.BindClaimME `json:"ys"`
}

type IdentME struct {
	DecID string `json:"id" param:"id"`
}

type DecRefME struct {
	DecID string `json:"dec_id" param:"id"`
	Title string `json:"title"`
	DecRN int64  `json:"dec_rn"`
}

type DecSnapME struct {
	X     termctx.BindClaimME   `json:"x"`
	DecID string                `json:"dec_id"`
	Ys    []termctx.BindClaimME `json:"ys"`
	Title string                `json:"title"`
	DecRN int64                 `json:"dec_rn"`
}
