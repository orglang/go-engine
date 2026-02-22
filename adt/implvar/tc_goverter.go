package implvar

import (
	"github.com/orglang/go-sdk/adt/implvar"
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
	DataFromVarRec  func(VarRec) VarRecDS
	DataFromVarRecs func([]VarRec) []VarRecDS
	// goverter:map . ImplRef
	DataToVarRec  func(VarRecDS) (VarRec, error)
	DataToVarRecs func([]VarRecDS) ([]VarRec, error)
)
