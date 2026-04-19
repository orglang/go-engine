package termexp

import (
	"github.com/orglang/go-sdk/pool/termexp"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend MsgFromExpSpec
// goverter:extend MsgToExpSpec
var (
	MsgFromAcquireSpec func(AcquireSpec) *termexp.AcquireSpec
	MsgFromAcceptSpec  func(AcceptSpec) *termexp.AcceptSpec
	MsgFromHireSpec    func(HireSpec) *termexp.HireSpec
	MsgFromApplySpec   func(ApplySpec) *termexp.ApplySpec
	MsgFromReleaseSpec func(ReleaseSpec) *termexp.ReleaseSpec
	MsgFromDetachSpec  func(DetachSpec) *termexp.DetachSpec

	MsgToAcquireSpec func(*termexp.AcquireSpec) (AcquireSpec, error)
	MsgToAcceptSpec  func(*termexp.AcceptSpec) (AcceptSpec, error)
	MsgToHireSpec    func(*termexp.HireSpec) (HireSpec, error)
	MsgToApplySpec   func(*termexp.ApplySpec) (ApplySpec, error)
	MsgToReleaseSpec func(*termexp.ReleaseSpec) (ReleaseSpec, error)
	MsgToDetachSpec  func(*termexp.DetachSpec) (DetachSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend DataFromExpSpec
// goverter:extend DataToExpSpec
var (
	DataFromAcquireSpec func(AcquireSpec) *upSpecDS
	DataFromAcceptSpec  func(AcceptSpec) *upSpecDS
	DataFromHireSpec    func(HireSpec) *laborSpecDS
	DataFromApplySpec   func(ApplySpec) *laborSpecDS
	DataFromReleaseSpec func(ReleaseSpec) *downSpecDS
	DataFromDetachSpec  func(DetachSpec) *downSpecDS

	DataToAcquireSpec func(*upSpecDS) (AcquireSpec, error)
	DataToAcceptSpec  func(*upSpecDS) (AcceptSpec, error)
	DataToHireSpec    func(*laborSpecDS) (HireSpec, error)
	DataToApplySpec   func(*laborSpecDS) (ApplySpec, error)
	DataToReleaseSpec func(*downSpecDS) (ReleaseSpec, error)
	DataToDetachSpec  func(*downSpecDS) (DetachSpec, error)

	DataFromAcquireRec func(AcquireRec) *upRecDS
	DataFromAcceptRec  func(AcceptRec) *upRecDS
	DataFromHireRec    func(HireRec) *laborRecDS
	DataFromApplyRec   func(ApplyRec) *laborRecDS
	DataFromReleaseRec func(ReleaseRec) *downRecDS
	DataFromDetachRec  func(DetachRec) *downRecDS

	DataToAcquireRec func(*upRecDS) (AcquireRec, error)
	DataToAcceptRec  func(*upRecDS) (AcceptRec, error)
	DataToHireRec    func(*laborRecDS) (HireRec, error)
	DataToApplyRec   func(*laborRecDS) (ApplyRec, error)
	DataToReleaseRec func(*downRecDS) (ReleaseRec, error)
	DataToDetachRec  func(*downRecDS) (DetachRec, error)
)
