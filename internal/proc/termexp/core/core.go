package core

import (
	"orglang/go-engine/internal/adt/identity"
	compexp "orglang/go-engine/internal/proc/compexec/exp"
)

type CompSpec = compexp.ExpSpec
type CompRec = compexp.ExpRec
type CloseSpec = compexp.CloseSpec
type WaitSpec = compexp.WaitSpec
type SendSpec = compexp.SendSpec
type RecvSpec = compexp.RecvSpec
type LabSpec = compexp.LabSpec
type CaseSpec = compexp.CaseSpec
type LinkSpec = compexp.LinkSpec
type FwdSpec = compexp.FwdSpec
type CallSpec = compexp.CallSpec
type SpawnSpec = compexp.SpawnSpec
type AcqureSpec = compexp.AcqureSpec
type AcceptSpec = compexp.AcceptSpec
type DetachSpec = compexp.DetachSpec
type ReleaseSpec = compexp.ReleaseSpec
type CloseRec = compexp.CloseRec
type WaitRec = compexp.WaitRec
type SendRec = compexp.SendRec
type RecvRec = compexp.RecvRec
type LabRec = compexp.LabRec
type CaseRec = compexp.CaseRec
type FwdRec = compexp.FwdRec

func CollectEnv(spec CompSpec) []identity.ADT { return compexp.CollectEnv(spec) }
func ErrExpTypeUnexpected(got CompSpec) error { return compexp.ErrExpTypeUnexpected(got) }
func ErrRecTypeUnexpected(got CompRec) error { return compexp.ErrRecTypeUnexpected(got) }
func ErrExpTypeMismatch(got, want CompSpec) error { return compexp.ErrExpTypeMismatch(got, want) }
func ErrExpValueNil(pid identity.ADT) error { return compexp.ErrExpValueNil(pid) }
