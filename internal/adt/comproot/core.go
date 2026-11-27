// Package comproot defines the Computation aggregate root reference.
package comproot

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

// RootRef identifies a Computation aggregate at a specific OCC revision.
// Execution and Connections are owned by this aggregate.
type RootRef struct {
	CompID identity.ADT
	CompRN seqnum.ADT
}

// New creates a reference for a new Computation aggregate.
func New() RootRef {
	return RootRef{CompID: identity.New(), CompRN: seqnum.New()}
}
