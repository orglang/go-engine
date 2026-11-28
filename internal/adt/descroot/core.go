// Package descroot defines the Description aggregate root reference.
package descroot

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

// RootRef identifies a Description aggregate at a specific OCC revision.
// The reference also guards changes to its owned Provider Type and Terms.
type RootRef struct {
	DescID identity.ADT
	DescRN seqnum.ADT
}

// New creates a reference for a new Description aggregate.
func New() RootRef {
	return RootRef{DescID: identity.New(), DescRN: seqnum.New()}
}
