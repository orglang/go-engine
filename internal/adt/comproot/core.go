// Package comproot defines the Computation aggregate root reference.
package comproot

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

// SemRef identifies a Computation aggregate at a specific OCC revision.
// Execution and Connections are owned by this aggregate.
type SemRef struct {
	CompID identity.ADT
	CompRN seqnum.ADT
}

// New creates a reference for a new Computation aggregate.
func New() SemRef {
	return SemRef{CompID: identity.New(), CompRN: seqnum.New()}
}
