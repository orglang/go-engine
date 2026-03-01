package poolcomm

import (
	"github.com/orglang/go-sdk/adt/poolcomm"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/implsem:Msg.*
// goverter:extend orglang/go-engine/adt/poolexp:Msg.*
var (
	MsgToCommSpec   func(poolcomm.CommSpec) (CommSpec, error)
	MsgFromCommSpec func(CommSpec) poolcomm.CommSpec
)
