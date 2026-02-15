package synonym

import (
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

type Rec struct {
	SynQN uniqsym.ADT
	SynVK valkey.ADT
}
