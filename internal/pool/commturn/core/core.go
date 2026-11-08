package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/commsem"
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

// step (aka Act)
type TurnRec interface {
	turn()
}

// publication (aka msg)
type PubRec struct {
	// совпадает со значением в SubRec
	CommRef commsem.SemRef
	CompRef compsem.SemRef
	ChnlID  identity.ADT
	ValExp  termexp.ExpRec
}

func (r PubRec) turn() {}

// subscription (aka srv)
type SubRec struct {
	// совпадает со значением в PubRec
	CommRef commsem.SemRef
	CompRef compsem.SemRef
	ChnlID  identity.ADT
	ContExp termexp.ExpRec
}

func (r SubRec) turn() {}

func ErrRecTypeUnexpected(got TurnRec) error {
	return fmt.Errorf("step rec unexpected: %T%+v", got, got)
}

func ErrRecTypeMismatch(got, want TurnRec) error {
	return fmt.Errorf("step rec mismatch: want %T, got %T", want, got)
}
