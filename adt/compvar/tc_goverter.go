package compvar

import (
	"github.com/orglang/go-sdk/adt/compvar"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/compsem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend ConvertRecToRef
var (
	ConvertRecsToRefs func([]LinearRec) []compsem.SemRef
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/symbol:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	MsgFromVarSpec  func(VarSpec) compvar.VarSpec
	MsgFromVarSpecs func([]VarSpec) []compvar.VarSpec
	MsgToVarSpec    func(compvar.VarSpec) (VarSpec, error)
	MsgToVarSpecs   func([]compvar.VarSpec) ([]VarSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
// goverter:extend orglang/go-engine/adt/symbol:Convert.*
// goverter:extend Convert.*
var (
	// goverter:autoMap CompRef
	// goverter:autoMap CommRef
	DataFromStructRec  func(StructRec) VarRecDS
	DataFromStructRecs func([]StructRec) []VarRecDS
	// goverter:autoMap CompRef
	// goverter:autoMap CommRef
	DataFromLinearRec  func(LinearRec) VarRecDS
	DataFromLinearRecs func([]LinearRec) []VarRecDS
	// goverter:map . CompRef
	// goverter:map . CommRef | dataToSemRef
	DataToStructRec  func(VarRecDS) (StructRec, error)
	DataToStructRecs func([]VarRecDS) ([]StructRec, error)
	// goverter:map . CompRef
	// goverter:map . CommRef | dataToSemRef
	DataToLinearRec  func(VarRecDS) (LinearRec, error)
	DataToLinearRecs func([]VarRecDS) ([]LinearRec, error)
	// goverter:ignore CommRN
	dataToSemRef func(VarRecDS) (commsem.SemRef, error)
)
