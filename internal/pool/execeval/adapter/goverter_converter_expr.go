package adapter

import (
	compexpr "orglang/go-engine/internal/pool/execeval/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend DataFromExprSpec
// goverter:extend DataToExprSpec
var (
	DataFromAcquireSpec func(compexpr.AcquireSpec) *GrantSpec
	DataFromAcceptSpec  func(compexpr.AcceptSpec) *GrantSpec
	DataFromHireSpec    func(compexpr.HireSpec) *CoopSpec
	DataFromApplySpec   func(compexpr.ApplySpec) *CoopSpec
	DataFromReleaseSpec func(compexpr.ReleaseSpec) *RevokeSpec
	DataFromDetachSpec  func(compexpr.DetachSpec) *RevokeSpec

	DataToAcquireSpec func(*GrantSpec) (compexpr.AcquireSpec, error)
	DataToAcceptSpec  func(*GrantSpec) (compexpr.AcceptSpec, error)
	DataToHireSpec    func(*CoopSpec) (compexpr.HireSpec, error)
	DataToApplySpec   func(*CoopSpec) (compexpr.ApplySpec, error)
	DataToReleaseSpec func(*RevokeSpec) (compexpr.ReleaseSpec, error)
	DataToDetachSpec  func(*RevokeSpec) (compexpr.DetachSpec, error)

	DataFromAcquireRec func(compexpr.AcquireRec) *GrantRec
	DataFromAcceptRec  func(compexpr.AcceptRec) *GrantRec
	DataFromHireRec    func(compexpr.HireRec) *CoopRec
	DataFromApplyRec   func(compexpr.ApplyRec) *CoopRec
	DataFromReleaseRec func(compexpr.ReleaseRec) *RevokeRec
	DataFromDetachRec  func(compexpr.DetachRec) *RevokeRec

	DataToAcquireRec func(*GrantRec) (compexpr.AcquireRec, error)
	DataToAcceptRec  func(*GrantRec) (compexpr.AcceptRec, error)
	DataToHireRec    func(*CoopRec) (compexpr.HireRec, error)
	DataToApplyRec   func(*CoopRec) (compexpr.ApplyRec, error)
	DataToReleaseRec func(*RevokeRec) (compexpr.ReleaseRec, error)
	DataToDetachRec  func(*RevokeRec) (compexpr.DetachRec, error)
)
