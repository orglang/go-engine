package adapter

import (
	execconn "orglang/go-engine/internal/pool/execconn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/identity/adapter:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum/adapter:Convert.*
// goverter:extend orglang/go-engine/internal/adt/execroot/adapter:Data.*
// goverter:extend orglang/go-engine/internal/pool/connact/adapter:Data.*
var (
	// goverter:map ConnID | identityToString
	DataFromQry func(execconn.ConnQry) ConnQry
	// goverter:map ConnID | identityToString
	// goverter:ignore ActNr
	DataFromMod func(execconn.ConnMod) ConnMod
	// goverter:map . ConnID | dataToConnID
	DataToSnap func(ConnSnap) (execconn.ConnSnap, error)
)
