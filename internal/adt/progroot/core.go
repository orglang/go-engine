// Package progroot defines the Program aggregate root reference.
package progroot

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

// RootRef identifies a Program aggregate at a specific OCC revision.
// The reference also guards changes to its owned Provider Type and Terms.
type RootRef struct {
	ProgID identity.ADT
	ProgRN seqnum.ADT
}

// New creates a reference for a new Program aggregate.
func New() RootRef {
	return RootRef{ProgID: identity.New(), ProgRN: seqnum.New()}
}
