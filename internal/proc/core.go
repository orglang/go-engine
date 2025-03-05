package proc

import (
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
	"smecalculus/rolevod/lib/rev"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/step"
)

// aka Configuration
type Snap struct {
	PoolID id.ADT
	ProcID id.ADT
	Chnls  map[ph.ADT]Chnl
	Steps  map[chnl.ID]step.Root
	Rev    rev.ADT
}

type Chnl struct {
	ChnlID  id.ADT
	StateID id.ADT
	// Provider Side
	PS EP
	// Client Side
	CS EP
}

type EP struct {
	PoolID id.ADT
	ProcID id.ADT
	Rev    rev.ADT
}

func ChnlPH(ch Chnl) ph.ADT { return ch.ChnlPH }

type Lock struct {
	PoolID id.ADT
	Rev    rev.ADT
}

type Mod struct {
	Bnds  []Bnd
	Steps []step.Root
	Locks []Lock
}

type Bnd struct {
	ProcID  id.ADT
	ChnlPH  ph.ADT
	ChnlID  id.ADT
	StateID id.ADT
	Rev     rev.ADT
}
