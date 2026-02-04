package uniqref

import (
	"github.com/orglang/go-sdk/adt/uniqref"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
var (
	MsgToADT    func(uniqref.Msg) (ADT, error)
	MsgFromADT  func(ADT) uniqref.Msg
	MsgToADTs   func([]uniqref.Msg) ([]ADT, error)
	MsgFromADTs func([]ADT) []uniqref.Msg
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
var (
	DataToADT    func(Data) (ADT, error)
	DataFromADT  func(ADT) Data
	DataToADTs   func([]Data) ([]ADT, error)
	DataFromADTs func([]ADT) []Data
)
