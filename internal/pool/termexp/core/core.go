package core

import (
	"fmt"
	compexp "orglang/go-engine/internal/pool/compexec/exp"
)

type ExpKind int

const (
	UnkKind ExpKind = iota
	AcquireKind
	AcceptKind
	HireKind
	ApplyKind
	ReleaseKind
	DetachKind
	FireKind
	QuitKind
	SpawnKind
)

type CompSpec = compexp.ExpSpec
type CompRec = compexp.ExpRec
type AcquireSpec = compexp.AcquireSpec
type AcceptSpec = compexp.AcceptSpec
type HireSpec = compexp.HireSpec
type ApplySpec = compexp.ApplySpec
type ReleaseSpec = compexp.ReleaseSpec
type DetachSpec = compexp.DetachSpec
type FireSpec = compexp.FireSpec
type QuitSpec = compexp.QuitSpec
type SpawnSpec = compexp.SpawnSpec
type SpawnSpec2 = compexp.SpawnSpec2
type AcquireRec = compexp.AcquireRec
type AcceptRec = compexp.AcceptRec
type HireRec = compexp.HireRec
type ApplyRec = compexp.ApplyRec
type ReleaseRec = compexp.ReleaseRec
type DetachRec = compexp.DetachRec
type FireRec = compexp.FireRec
type QuitRec = compexp.QuitRec
type SpawnRec = compexp.SpawnRec

func ErrSpecTypeUnexpected(got CompSpec) error { return compexp.ErrSpecTypeUnexpected(got) }
func ErrRecTypeUnexpected(got CompRec) error { return compexp.ErrRecTypeUnexpected(got) }
func ErrExpKindUnexpected(got ExpKind) error { return fmt.Errorf("exp kind unexpected: %v", got) }
