package core

import termdef "orglang/go-engine/internal/pool/termdef/core"

type DescSpec = termdef.ExpSpec
type DescAcquireSpec = termdef.ExpAcquireSpec
type DescAcceptSpec = termdef.ExpAcceptSpec
type DescHireSpec = termdef.ExpHireSpec
type DescApplySpec = termdef.ExpApplySpec
type DescReleaseSpec = termdef.ExpReleaseSpec
type DescDetachSpec = termdef.ExpDetachSpec
type DescFireSpec = termdef.ExpFireSpec
type DescQuitSpec = termdef.ExpQuitSpec
type DescSpawnSpec = termdef.ExpSpawnSpec
type DescRec = termdef.ExpRec
type DescAcquireRec = termdef.ExpAcquireRec
type DescAcceptRec = termdef.ExpAcceptRec
type DescHireRec = termdef.ExpHireRec
type DescApplyRec = termdef.ExpApplyRec
type DescReleaseRec = termdef.ExpReleaseRec
type DescDetachRec = termdef.ExpDetachRec
type DescFireRec = termdef.ExpFireRec
type DescQuitRec = termdef.ExpQuitRec
type DescSpawnRec = termdef.ExpSpawnRec

func ConvertSpecToRec(spec DescSpec) (DescRec, error) { return termdef.ConvertSpecToRec(spec) }
func ConvertRecToSpec(rec DescRec) DescSpec { return termdef.ConvertRecToSpec(rec) }
