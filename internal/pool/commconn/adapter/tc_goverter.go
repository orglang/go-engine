package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	commconn "orglang/go-engine/internal/pool/commconn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compsem:Data.*
// goverter:extend orglang/go-engine/internal/pool/commturn/adapter:Data.*
var (
	// goverter:map ConnID | identityToString
	// goverter:autoMap CompRef
	DataFromRec func(commconn.ConnRec) commconn.ConnRecDS
	// goverter:map ConnID | identityToString
	DataFromQry func(commconn.ConnQry) commconn.ConnQryDS
	// goverter:map ConnID | identityToString
	// goverter:ignore Turns
	DataFromMod func(commconn.ConnMod) commconn.ConnModDS
	// goverter:ignore CompID
	// goverter:ignore CompRN
	// goverter:ignore ActNr
	// goverter:map . ConnID | dataToConnID
	DataToSnap func(commconn.ConnSnapDS) (commconn.ConnSnap, error)
)
