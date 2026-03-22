package implvar

import (
	"github.com/orglang/go-sdk/adt/implvar"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/implsem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend ConvertRecToRef
var (
	ConvertRecsToRefs func([]LinearRec) []implsem.SemRef
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/symbol:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	MsgFromVarSpec  func(VarSpec) implvar.VarSpec
	MsgFromVarSpecs func([]VarSpec) []implvar.VarSpec
	MsgToVarSpec    func(implvar.VarSpec) (VarSpec, error)
	MsgToVarSpecs   func([]implvar.VarSpec) ([]VarSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
var (
	// goverter:autoMap ImplRef
	// goverter:autoMap CommRef
	DataFromStructRec  func(StructRec) VarRecDS
	DataFromStructRecs func([]StructRec) []VarRecDS
	// goverter:autoMap ImplRef
	// goverter:autoMap CommRef
	DataFromLinearRec  func(LinearRec) VarRecDS
	DataFromLinearRecs func([]LinearRec) []VarRecDS
	// goverter:map . ImplRef
	// goverter:map . CommRef | DataToCommRef
	DataToStructRec  func(VarRecDS) (StructRec, error)
	DataToStructRecs func([]VarRecDS) ([]StructRec, error)
	// goverter:map . ImplRef
	// goverter:map . CommRef | DataToCommRef
	DataToLinearRec  func(VarRecDS) (LinearRec, error)
	DataToLinearRecs func([]VarRecDS) ([]LinearRec, error)
	// goverter:ignore CommRN
	DataToCommRef func(VarRecDS) (commsem.SemRef, error)
)
