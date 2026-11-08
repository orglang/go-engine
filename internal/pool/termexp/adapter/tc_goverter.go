package adapter

import (
	sdktermexp "github.com/orglang/go-sdk/pool/termexp"

	termexp "orglang/go-engine/pool/termexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend MsgFromExpSpec
// goverter:extend MsgToExpSpec
var (
	MsgFromAcquireSpec func(termexp.AcquireSpec) *sdktermexp.AcquireSpec
	MsgFromAcceptSpec  func(termexp.AcceptSpec) *sdktermexp.AcceptSpec
	MsgFromHireSpec    func(termexp.HireSpec) *sdktermexp.HireSpec
	MsgFromApplySpec   func(termexp.ApplySpec) *sdktermexp.ApplySpec
	MsgFromReleaseSpec func(termexp.ReleaseSpec) *sdktermexp.ReleaseSpec
	MsgFromDetachSpec  func(termexp.DetachSpec) *sdktermexp.DetachSpec

	MsgToAcquireSpec func(*sdktermexp.AcquireSpec) (termexp.AcquireSpec, error)
	MsgToAcceptSpec  func(*sdktermexp.AcceptSpec) (termexp.AcceptSpec, error)
	MsgToHireSpec    func(*sdktermexp.HireSpec) (termexp.HireSpec, error)
	MsgToApplySpec   func(*sdktermexp.ApplySpec) (termexp.ApplySpec, error)
	MsgToReleaseSpec func(*sdktermexp.ReleaseSpec) (termexp.ReleaseSpec, error)
	MsgToDetachSpec  func(*sdktermexp.DetachSpec) (termexp.DetachSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend DataFromExpSpec
// goverter:extend DataToExpSpec
var (
	DataFromAcquireSpec func(termexp.AcquireSpec) *termexp.GrantSpecDS
	DataFromAcceptSpec  func(termexp.AcceptSpec) *termexp.GrantSpecDS
	DataFromHireSpec    func(termexp.HireSpec) *termexp.CoopSpecDS
	DataFromApplySpec   func(termexp.ApplySpec) *termexp.CoopSpecDS
	DataFromReleaseSpec func(termexp.ReleaseSpec) *termexp.RevokeSpecDS
	DataFromDetachSpec  func(termexp.DetachSpec) *termexp.RevokeSpecDS

	DataToAcquireSpec func(*termexp.GrantSpecDS) (termexp.AcquireSpec, error)
	DataToAcceptSpec  func(*termexp.GrantSpecDS) (termexp.AcceptSpec, error)
	DataToHireSpec    func(*termexp.CoopSpecDS) (termexp.HireSpec, error)
	DataToApplySpec   func(*termexp.CoopSpecDS) (termexp.ApplySpec, error)
	DataToReleaseSpec func(*termexp.RevokeSpecDS) (termexp.ReleaseSpec, error)
	DataToDetachSpec  func(*termexp.RevokeSpecDS) (termexp.DetachSpec, error)

	DataFromAcquireRec func(termexp.AcquireRec) *termexp.GrantRecDS
	DataFromAcceptRec  func(termexp.AcceptRec) *termexp.GrantRecDS
	DataFromHireRec    func(termexp.HireRec) *termexp.CoopRecDS
	DataFromApplyRec   func(termexp.ApplyRec) *termexp.CoopRecDS
	DataFromReleaseRec func(termexp.ReleaseRec) *termexp.RevokeRecDS
	DataFromDetachRec  func(termexp.DetachRec) *termexp.RevokeRecDS

	DataToAcquireRec func(*termexp.GrantRecDS) (termexp.AcquireRec, error)
	DataToAcceptRec  func(*termexp.GrantRecDS) (termexp.AcceptRec, error)
	DataToHireRec    func(*termexp.CoopRecDS) (termexp.HireRec, error)
	DataToApplyRec   func(*termexp.CoopRecDS) (termexp.ApplyRec, error)
	DataToReleaseRec func(*termexp.RevokeRecDS) (termexp.ReleaseRec, error)
	DataToDetachRec  func(*termexp.RevokeRecDS) (termexp.DetachRec, error)
)
