package adapter

import (
	conncore "orglang/go-engine/internal/pool/conn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compsem:Data.*
// goverter:extend orglang/go-engine/internal/pool/connact/adapter:Data.*
var (
	// goverter:map ConnID | identityToString
	DataFromQry func(conncore.ConnQry) conncore.ConnQryDS
	// goverter:map ConnID | identityToString
	// goverter:ignore ActNr
	DataFromMod func(conncore.ConnMod) conncore.ConnModDS
	// goverter:map . ConnID | dataToConnID
	DataToSnap func(conncore.ConnSnapDS) (conncore.ConnSnap, error)
)
