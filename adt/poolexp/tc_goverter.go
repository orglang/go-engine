package poolexp

import (
	"github.com/orglang/go-sdk/adt/poolexp"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend MsgFromExpSpec
// goverter:extend MsgToExpSpec
var (
	MsgFromAcquireSpec func(AcquireSpec) *poolexp.AcquireSpec
	MsgFromAcceptSpec  func(AcceptSpec) *poolexp.AcceptSpec
	MsgFromHireSpec    func(HireSpec) *poolexp.HireSpec
	MsgFromApplySpec   func(ApplySpec) *poolexp.ApplySpec
	MsgFromReleaseSpec func(ReleaseSpec) *poolexp.ReleaseSpec
	MsgFromDetachSpec  func(DetachSpec) *poolexp.DetachSpec

	MsgToAcquireSpec func(*poolexp.AcquireSpec) (AcquireSpec, error)
	MsgToAcceptSpec  func(*poolexp.AcceptSpec) (AcceptSpec, error)
	MsgToHireSpec    func(*poolexp.HireSpec) (HireSpec, error)
	MsgToApplySpec   func(*poolexp.ApplySpec) (ApplySpec, error)
	MsgToReleaseSpec func(*poolexp.ReleaseSpec) (ReleaseSpec, error)
	MsgToDetachSpec  func(*poolexp.DetachSpec) (DetachSpec, error)
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
