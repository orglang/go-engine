// Package descroot defines the Description aggregate root reference.
package descroot

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

// SemRef identifies a Description aggregate at a specific OCC revision.
// The reference also guards changes to its owned Provider Type and Terms.
type SemRef struct {
	DescID identity.ADT
	DescRN seqnum.ADT
}

// New creates a reference for a new Description aggregate.
func New() SemRef {
	return SemRef{DescID: identity.New(), DescRN: seqnum.New()}
}
