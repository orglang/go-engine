package implvar

import (
	"github.com/orglang/go-sdk/adt/implvar"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/symbol:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	MsgToVarSpec   func(implvar.VarSpec) (VarSpec, error)
	MsgFromVarSpec func(VarSpec) implvar.VarSpec
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
var (
	DataToVarSpec   func(VarSpecDS) (VarSpec, error)
	DataFromVarSpec func(VarSpec) VarSpecDS
	// goverter:map . ImplRef
	DataToVarRec func(VarRecDS) (VarRec, error)
	// goverter:autoMap ImplRef
	DataFromVarRec  func(VarRec) VarRecDS
	DataToVarRecs   func([]VarRecDS) ([]VarRec, error)
	DataFromVarRecs func([]VarRec) []VarRecDS
)
