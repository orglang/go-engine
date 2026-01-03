package syndec

import (
	"orglang/orglang/lib/sd"
)

type Repo interface {
	Insert(sd.Source, DecRec) error
}

type decRecDS struct {
	DecID string
	DecRN int64
	DecQN string
}
