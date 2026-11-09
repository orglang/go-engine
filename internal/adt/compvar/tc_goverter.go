package compvar

import (
	"database/sql"

	"github.com/orglang/go-sdk/adt/compvar"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend ConvertRecToRef
var (
	ConvertRecsToRefs func([]LinearRec) []compsem.SemRef
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/symbol:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
var (
	MsgFromVarSpec  func(VarSpec) compvar.VarSpec
	MsgFromVarSpecs func([]VarSpec) []compvar.VarSpec
	MsgToVarSpec    func(compvar.VarSpec) (VarSpec, error)
	MsgToVarSpecs   func([]compvar.VarSpec) ([]VarSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/valkey:Convert.*
// goverter:extend orglang/go-engine/internal/adt/symbol:Convert.*
// goverter:extend Convert.*
var (
	// goverter:autoMap CompRef
	// goverter:map ConnID | identityToNullString
	DataFromStructRec  func(StructRec) VarRecDS
	DataFromStructRecs func([]StructRec) []VarRecDS
	// goverter:autoMap CompRef
	// goverter:map ConnID | identityToNullString
	DataFromLinearRec  func(LinearRec) VarRecDS
	DataFromLinearRecs func([]LinearRec) []VarRecDS
	// goverter:map . CompRef
	// goverter:map ConnID | dataToConnID
	DataToStructRec  func(VarRecDS) (StructRec, error)
	DataToStructRecs func([]VarRecDS) ([]StructRec, error)
	// goverter:map . CompRef
	// goverter:map ConnID | dataToConnID
	DataToLinearRec  func(VarRecDS) (LinearRec, error)
	DataToLinearRecs func([]VarRecDS) ([]LinearRec, error)
)

func identityToNullString(id identity.ADT) sql.NullString {
	return identity.ConvertToNullString(id)
}

func dataToConnID(dto sql.NullString) (identity.ADT, error) {
	return identity.ConvertFromNullString(dto)
}
