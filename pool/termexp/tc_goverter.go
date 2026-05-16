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
	DataFromAcquireSpec func(AcquireSpec) *grantSpecDS
	DataFromAcceptSpec  func(AcceptSpec) *grantSpecDS
	DataFromHireSpec    func(HireSpec) *coopSpecDS
	DataFromApplySpec   func(ApplySpec) *coopSpecDS
	DataFromReleaseSpec func(ReleaseSpec) *revokeSpecDS
	DataFromDetachSpec  func(DetachSpec) *revokeSpecDS

	DataToAcquireSpec func(*grantSpecDS) (AcquireSpec, error)
	DataToAcceptSpec  func(*grantSpecDS) (AcceptSpec, error)
	DataToHireSpec    func(*coopSpecDS) (HireSpec, error)
	DataToApplySpec   func(*coopSpecDS) (ApplySpec, error)
	DataToReleaseSpec func(*revokeSpecDS) (ReleaseSpec, error)
	DataToDetachSpec  func(*revokeSpecDS) (DetachSpec, error)

	DataFromAcquireRec func(AcquireRec) *grantRecDS
	DataFromAcceptRec  func(AcceptRec) *grantRecDS
	DataFromHireRec    func(HireRec) *coopRecDS
	DataFromApplyRec   func(ApplyRec) *coopRecDS
	DataFromReleaseRec func(ReleaseRec) *revokeRecDS
	DataFromDetachRec  func(DetachRec) *revokeRecDS

	DataToAcquireRec func(*grantRecDS) (AcquireRec, error)
	DataToAcceptRec  func(*grantRecDS) (AcceptRec, error)
	DataToHireRec    func(*coopRecDS) (HireRec, error)
	DataToApplyRec   func(*coopRecDS) (ApplyRec, error)
	DataToReleaseRec func(*revokeRecDS) (ReleaseRec, error)
	DataToDetachRec  func(*revokeRecDS) (DetachRec, error)
)
