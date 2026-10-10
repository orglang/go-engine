package adapter

import (
	"github.com/orglang/go-sdk/pool/termexp"

	pooltermexp "orglang/go-engine/pool/termexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend MsgFromExpSpec
// goverter:extend MsgToExpSpec
var (
	MsgFromAcquireSpec func(pooltermexp.AcquireSpec) *termexp.AcquireSpec
	MsgFromAcceptSpec  func(pooltermexp.AcceptSpec) *termexp.AcceptSpec
	MsgFromHireSpec    func(pooltermexp.HireSpec) *termexp.HireSpec
	MsgFromApplySpec   func(pooltermexp.ApplySpec) *termexp.ApplySpec
	MsgFromReleaseSpec func(pooltermexp.ReleaseSpec) *termexp.ReleaseSpec
	MsgFromDetachSpec  func(pooltermexp.DetachSpec) *termexp.DetachSpec

	MsgToAcquireSpec func(*termexp.AcquireSpec) (pooltermexp.AcquireSpec, error)
	MsgToAcceptSpec  func(*termexp.AcceptSpec) (pooltermexp.AcceptSpec, error)
	MsgToHireSpec    func(*termexp.HireSpec) (pooltermexp.HireSpec, error)
	MsgToApplySpec   func(*termexp.ApplySpec) (pooltermexp.ApplySpec, error)
	MsgToReleaseSpec func(*termexp.ReleaseSpec) (pooltermexp.ReleaseSpec, error)
	MsgToDetachSpec  func(*termexp.DetachSpec) (pooltermexp.DetachSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend DataFromExpSpec
// goverter:extend DataToExpSpec
var (
	DataFromAcquireSpec func(pooltermexp.AcquireSpec) *pooltermexp.GrantSpecDS
	DataFromAcceptSpec  func(pooltermexp.AcceptSpec) *pooltermexp.GrantSpecDS
	DataFromHireSpec    func(pooltermexp.HireSpec) *pooltermexp.CoopSpecDS
	DataFromApplySpec   func(pooltermexp.ApplySpec) *pooltermexp.CoopSpecDS
	DataFromReleaseSpec func(pooltermexp.ReleaseSpec) *pooltermexp.RevokeSpecDS
	DataFromDetachSpec  func(pooltermexp.DetachSpec) *pooltermexp.RevokeSpecDS

	DataToAcquireSpec func(*pooltermexp.GrantSpecDS) (pooltermexp.AcquireSpec, error)
	DataToAcceptSpec  func(*pooltermexp.GrantSpecDS) (pooltermexp.AcceptSpec, error)
	DataToHireSpec    func(*pooltermexp.CoopSpecDS) (pooltermexp.HireSpec, error)
	DataToApplySpec   func(*pooltermexp.CoopSpecDS) (pooltermexp.ApplySpec, error)
	DataToReleaseSpec func(*pooltermexp.RevokeSpecDS) (pooltermexp.ReleaseSpec, error)
	DataToDetachSpec  func(*pooltermexp.RevokeSpecDS) (pooltermexp.DetachSpec, error)

	DataFromAcquireRec func(pooltermexp.AcquireRec) *pooltermexp.GrantRecDS
	DataFromAcceptRec  func(pooltermexp.AcceptRec) *pooltermexp.GrantRecDS
	DataFromHireRec    func(pooltermexp.HireRec) *pooltermexp.CoopRecDS
	DataFromApplyRec   func(pooltermexp.ApplyRec) *pooltermexp.CoopRecDS
	DataFromReleaseRec func(pooltermexp.ReleaseRec) *pooltermexp.RevokeRecDS
	DataFromDetachRec  func(pooltermexp.DetachRec) *pooltermexp.RevokeRecDS

	DataToAcquireRec func(*pooltermexp.GrantRecDS) (pooltermexp.AcquireRec, error)
	DataToAcceptRec  func(*pooltermexp.GrantRecDS) (pooltermexp.AcceptRec, error)
	DataToHireRec    func(*pooltermexp.CoopRecDS) (pooltermexp.HireRec, error)
	DataToApplyRec   func(*pooltermexp.CoopRecDS) (pooltermexp.ApplyRec, error)
	DataToReleaseRec func(*pooltermexp.RevokeRecDS) (pooltermexp.ReleaseRec, error)
	DataToDetachRec  func(*pooltermexp.RevokeRecDS) (pooltermexp.DetachRec, error)
)
