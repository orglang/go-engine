package adapter

import (
	sdktermexp "github.com/orglang/go-sdk/pool/termexp"

	termexp "orglang/go-engine/internal/pool/termexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend MsgFromExpSpec
// goverter:extend MsgToExpSpec
var (
	// goverter:map ConnChnlPH CommChnlPH
	MsgFromAcquireSpec func(termexp.AcquireSpec) *sdktermexp.AcquireSpec
	// goverter:map ConnChnlPH CommChnlPH
	MsgFromAcceptSpec  func(termexp.AcceptSpec) *sdktermexp.AcceptSpec
	// goverter:map ConnChnlPH CommChnlPH
	MsgFromHireSpec    func(termexp.HireSpec) *sdktermexp.HireSpec
	// goverter:map ConnChnlPH CommChnlPH
	MsgFromApplySpec   func(termexp.ApplySpec) *sdktermexp.ApplySpec
	// goverter:map ConnChnlPH CommChnlPH
	MsgFromReleaseSpec func(termexp.ReleaseSpec) *sdktermexp.ReleaseSpec
	// goverter:map ConnChnlPH CommChnlPH
	MsgFromDetachSpec  func(termexp.DetachSpec) *sdktermexp.DetachSpec

	// goverter:map CommChnlPH ConnChnlPH
	MsgToAcquireSpec func(*sdktermexp.AcquireSpec) (termexp.AcquireSpec, error)
	// goverter:map CommChnlPH ConnChnlPH
	MsgToAcceptSpec  func(*sdktermexp.AcceptSpec) (termexp.AcceptSpec, error)
	// goverter:map CommChnlPH ConnChnlPH
	MsgToHireSpec    func(*sdktermexp.HireSpec) (termexp.HireSpec, error)
	// goverter:map CommChnlPH ConnChnlPH
	MsgToApplySpec   func(*sdktermexp.ApplySpec) (termexp.ApplySpec, error)
	// goverter:map CommChnlPH ConnChnlPH
	MsgToReleaseSpec func(*sdktermexp.ReleaseSpec) (termexp.ReleaseSpec, error)
	// goverter:map CommChnlPH ConnChnlPH
	MsgToDetachSpec  func(*sdktermexp.DetachSpec) (termexp.DetachSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
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
