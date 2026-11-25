package adapter

import (
	compexp "orglang/go-engine/internal/pool/compexec/exp"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend DataFromCompSpec
// goverter:extend DataToCompSpec
var (
	DataFromAcquireSpec func(termexp.AcquireSpec) *GrantSpec
	DataFromAcceptSpec  func(termexp.AcceptSpec) *GrantSpec
	DataFromHireSpec    func(termexp.HireSpec) *CoopSpec
	DataFromApplySpec   func(termexp.ApplySpec) *CoopSpec
	DataFromReleaseSpec func(termexp.ReleaseSpec) *RevokeSpec
	DataFromDetachSpec  func(termexp.DetachSpec) *RevokeSpec

	DataToAcquireSpec func(*GrantSpec) (termexp.AcquireSpec, error)
	DataToAcceptSpec  func(*GrantSpec) (termexp.AcceptSpec, error)
	DataToHireSpec    func(*CoopSpec) (termexp.HireSpec, error)
	DataToApplySpec   func(*CoopSpec) (termexp.ApplySpec, error)
	DataToReleaseSpec func(*RevokeSpec) (termexp.ReleaseSpec, error)
	DataToDetachSpec  func(*RevokeSpec) (termexp.DetachSpec, error)

	DataFromAcquireRec func(termexp.AcquireRec) *GrantRec
	DataFromAcceptRec  func(termexp.AcceptRec) *GrantRec
	DataFromHireRec    func(termexp.HireRec) *CoopRec
	DataFromApplyRec   func(termexp.ApplyRec) *CoopRec
	DataFromReleaseRec func(termexp.ReleaseRec) *RevokeRec
	DataFromDetachRec  func(termexp.DetachRec) *RevokeRec

	DataToAcquireRec func(*GrantRec) (termexp.AcquireRec, error)
	DataToAcceptRec  func(*GrantRec) (termexp.AcceptRec, error)
	DataToHireRec    func(*CoopRec) (termexp.HireRec, error)
	DataToApplyRec   func(*CoopRec) (termexp.ApplyRec, error)
	DataToReleaseRec func(*RevokeRec) (termexp.ReleaseRec, error)
	DataToDetachRec  func(*RevokeRec) (termexp.DetachRec, error)
)
