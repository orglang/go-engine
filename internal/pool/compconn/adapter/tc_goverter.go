package adapter

import (
	compconncore "orglang/go-engine/internal/pool/compconn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compsem:Data.*
// goverter:extend orglang/go-engine/internal/pool/connact/adapter:Data.*
var (
	// goverter:map ConnID | identityToString
	DataFromQry func(compconncore.ConnQry) compconncore.ConnQryDS
	// goverter:map ConnID | identityToString
	// goverter:ignore ActNr
	DataFromMod func(compconncore.ConnMod) compconncore.ConnModDS
	// goverter:map . ConnID | dataToConnID
	DataToSnap func(compconncore.ConnSnapDS) (compconncore.ConnSnap, error)
)
