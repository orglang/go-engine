package adapter

import (
	execcomm "orglang/go-engine/internal/pool/execcomm/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/identity/adapter:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum/adapter:Convert.*
// goverter:extend orglang/go-engine/internal/adt/execroot/adapter:Data.*
// goverter:extend orglang/go-engine/internal/pool/commact/adapter:Data.*
var (
	// goverter:map CommID | identityToString
	DataFromQry func(execcomm.CommQry) CommQry
	// goverter:map CommID | identityToString
	// goverter:ignore ActNr
	DataFromMod func(execcomm.CommMod) CommMod
	// goverter:map . CommID | dataToCommID
	DataToSnap func(CommSnap) (execcomm.CommSnap, error)
)
