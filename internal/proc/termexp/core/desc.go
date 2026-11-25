package core

import termdef "orglang/go-engine/internal/proc/termdef/core"

type DescSpec = termdef.ExpSpec
type DescCloseSpec = termdef.ExpCloseSpec
type DescWaitSpec = termdef.ExpWaitSpec
type DescSendSpec = termdef.ExpSendSpec
type DescRecvSpec = termdef.ExpRecvSpec
type DescLabSpec = termdef.ExpLabSpec
type DescCaseSpec = termdef.ExpCaseSpec
type DescLinkSpec = termdef.ExpLinkSpec
type DescFwdSpec = termdef.ExpFwdSpec
type DescCallSpec = termdef.ExpCallSpec
type DescSpawnSpec = termdef.ExpSpawnSpec
type DescAcquireSpec = termdef.ExpAcquireSpec
type DescAcceptSpec = termdef.ExpAcceptSpec
type DescDetachSpec = termdef.ExpDetachSpec
type DescReleaseSpec = termdef.ExpReleaseSpec
type DescRec = termdef.ExpRec
type DescCloseRec = termdef.ExpCloseRec
type DescWaitRec = termdef.ExpWaitRec
type DescSendRec = termdef.ExpSendRec
type DescRecvRec = termdef.ExpRecvRec
type DescLabRec = termdef.ExpLabRec
type DescCaseRec = termdef.ExpCaseRec
type DescLinkRec = termdef.ExpLinkRec
type DescFwdRec = termdef.ExpFwdRec
type DescCallRec = termdef.ExpCallRec
type DescSpawnRec = termdef.ExpSpawnRec
type DescAcquireRec = termdef.ExpAcquireRec
type DescAcceptRec = termdef.ExpAcceptRec
type DescDetachRec = termdef.ExpDetachRec
type DescReleaseRec = termdef.ExpReleaseRec

func ConvertSpecToRec(spec DescSpec) (DescRec, error) { return termdef.ConvertSpecToRec(spec) }
func ConvertRecToSpec(rec DescRec) DescSpec { return termdef.ConvertRecToSpec(rec) }
