package adapter

import (
	compexp "orglang/go-engine/internal/pool/compexec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend DataFromExpSpec
// goverter:extend DataToExpSpec
var (
	DataFromAcquireSpec func(compexp.AcquireSpec) *GrantSpec
	DataFromAcceptSpec  func(compexp.AcceptSpec) *GrantSpec
	DataFromHireSpec    func(compexp.HireSpec) *CoopSpec
	DataFromApplySpec   func(compexp.ApplySpec) *CoopSpec
	DataFromReleaseSpec func(compexp.ReleaseSpec) *RevokeSpec
	DataFromDetachSpec  func(compexp.DetachSpec) *RevokeSpec

	DataToAcquireSpec func(*GrantSpec) (compexp.AcquireSpec, error)
	DataToAcceptSpec  func(*GrantSpec) (compexp.AcceptSpec, error)
	DataToHireSpec    func(*CoopSpec) (compexp.HireSpec, error)
	DataToApplySpec   func(*CoopSpec) (compexp.ApplySpec, error)
	DataToReleaseSpec func(*RevokeSpec) (compexp.ReleaseSpec, error)
	DataToDetachSpec  func(*RevokeSpec) (compexp.DetachSpec, error)

	DataFromAcquireRec func(compexp.AcquireRec) *GrantRec
	DataFromAcceptRec  func(compexp.AcceptRec) *GrantRec
	DataFromHireRec    func(compexp.HireRec) *CoopRec
	DataFromApplyRec   func(compexp.ApplyRec) *CoopRec
	DataFromReleaseRec func(compexp.ReleaseRec) *RevokeRec
	DataFromDetachRec  func(compexp.DetachRec) *RevokeRec

	DataToAcquireRec func(*GrantRec) (compexp.AcquireRec, error)
	DataToAcceptRec  func(*GrantRec) (compexp.AcceptRec, error)
	DataToHireRec    func(*CoopRec) (compexp.HireRec, error)
	DataToApplyRec   func(*CoopRec) (compexp.ApplyRec, error)
	DataToReleaseRec func(*RevokeRec) (compexp.ReleaseRec, error)
	DataToDetachRec  func(*RevokeRec) (compexp.DetachRec, error)
)
