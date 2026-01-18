package procdec

import (
	"github.com/orglang/go-sdk/adt/procdec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-runtime/adt/identity:Convert.*
var (
	ConvertSnapToRef func(DecSnap) DecRef
	ConvertRecToRef  func(DecRec) DecRef
	ConvertRecToSnap func(DecRec) DecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-runtime/adt/uniqsym:Convert.*
// goverter:extend orglang/go-runtime/adt/termctx:Msg.*
// goverter:extend orglang/go-runtime/adt/typedef:Msg.*
var (
	MsgToDecSpec    func(procdec.DecSpec) (DecSpec, error)
	MsgFromDecSpec  func(DecSpec) procdec.DecSpec
	MsgToDecSnap    func(procdec.DecSnap) (DecSnap, error)
	MsgFromDecSnap  func(DecSnap) procdec.DecSnap
	MsgFromDecSnaps func([]DecSnap) []procdec.DecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-runtime/adt/identity:Convert.*
// goverter:extend orglang/go-runtime/adt/revnum:Convert.*
// goverter:extend orglang/go-runtime/adt/typedef:Msg.*
// goverter:extend Msg.*
var (
	ViewFromDecSnap func(DecSnap) DecSnapVP
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-runtime/adt/identity:Convert.*
// goverter:extend orglang/go-runtime/adt/termctx:Data.*
// goverter:extend orglang/go-runtime/adt/typedef:Data.*
var (
	DataToDecRec     func(decRecDS) (DecRec, error)
	DataFromDecRec   func(DecRec) (decRecDS, error)
	DataToDecRecs    func([]decRecDS) ([]DecRec, error)
	DataFromDecRecs  func([]DecRec) ([]decRecDS, error)
	DataToDecSnap    func(decSnapDS) (DecSnap, error)
	DataFromDecSnap  func(DecSnap) (decSnapDS, error)
	DataToDecSnaps   func([]decSnapDS) ([]DecSnap, error)
	DataFromDecSnaps func([]DecSnap) ([]decSnapDS, error)
)
