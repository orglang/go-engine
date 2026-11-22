// Package adt contains reusable domain value types shared across bounded contexts.
package adt

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

// DescRef identifies a Description aggregate and its optimistic-concurrency revision.
// The same reference also guards mutations of the Description's owned Provider Type
// and Terms; those entities do not have independent aggregate revisions.
type DescRef struct {
	DescID identity.ADT
	DescRN seqnum.ADT
}

// NewDescRef creates a reference to a new Description aggregate.
func NewDescRef() DescRef {
	return DescRef{DescID: identity.New(), DescRN: seqnum.New()}
}

// CompRef identifies a Computation aggregate and its independent OCC revision.
// Execution and Connections belong to this aggregate.
type CompRef struct {
	CompID identity.ADT
	CompRN seqnum.ADT
}

// NewCompRef creates a reference to a new Computation aggregate.
func NewCompRef() CompRef {
	return CompRef{CompID: identity.New(), CompRN: seqnum.New()}
}
