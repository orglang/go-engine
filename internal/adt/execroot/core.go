// Package execroot defines the execution aggregate root reference.
package execroot

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

// RootRef identifies an execution aggregate at a specific OCC revision.
// Executions and their connections are owned by this aggregate.
type RootRef struct {
	ExecID identity.ADT
	ExecRN seqnum.ADT
}

// New creates a reference for a new execution aggregate.
func New() RootRef {
	return RootRef{ExecID: identity.New(), ExecRN: seqnum.New()}
}
