package semcomp

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/semterm"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	MsgFromRef  func(CompRef) semterm.TermRef
	MsgFromRefs func([]CompRef) []semterm.TermRef
	MsgToRef    func(semterm.TermRef) (CompRef, error)
	MsgToRefs   func([]semterm.TermRef) ([]CompRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	DataFromRef  func(CompRef) CompRefDS
	DataFromRefs func([]CompRef) ([]CompRefDS, error)
	DataToRef    func(CompRefDS) (CompRef, error)
	DataToRefs   func([]CompRefDS) ([]CompRef, error)
	DataToRefMap func(map[uniqsym.ADT]CompRefDS) (map[uniqsym.ADT]CompRef, error)
)
