package adapter

import (
	commconn "orglang/go-engine/internal/pool/commconn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compsem:Data.*
// goverter:extend orglang/go-engine/internal/pool/commact/adapter:Data.*
var (
	// goverter:map ConnID | identityToString
	DataFromQry func(commconn.ConnQry) commconn.ConnQryDS
	// goverter:map ConnID | identityToString
	// goverter:ignore ActNr
	DataFromMod func(commconn.ConnMod) commconn.ConnModDS
	// goverter:map . ConnID | dataToConnID
	DataToSnap func(commconn.ConnSnapDS) (commconn.ConnSnap, error)
)
