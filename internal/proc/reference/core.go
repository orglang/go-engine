// Package reference defines the identity and optimistic-concurrency references
// used inside the proc bounded context.
//
// These types intentionally live in proc rather than adt: a proc reference
// cannot be passed to pool, even though both contexts use the same ID generator.
package reference

import (
	"context"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/uniqsym"
)

// Description, Type, Term and Computation identities are distinct types.
type (
	DescID identity.ADT
	TypeID identity.ADT
	TermID identity.ADT
	CompID identity.ADT
)

// Description and Computation are the only aggregate roots in proc.
type (
	DescRN seqnum.ADT
	CompRN seqnum.ADT
)

func NewDescID() DescID { return DescID(identity.New()) }
func NewTypeID() TypeID { return TypeID(identity.New()) }
func NewTermID() TermID { return TermID(identity.New()) }
func NewCompID() CompID { return CompID(identity.New()) }

func NewDescRN() DescRN { return DescRN(seqnum.New()) }
func NewCompRN() CompRN { return CompRN(seqnum.New()) }

func (id DescID) IsEmpty() bool { return identity.ADT(id).IsEmpty() }
func (id TypeID) IsEmpty() bool { return identity.ADT(id).IsEmpty() }
func (id TermID) IsEmpty() bool { return identity.ADT(id).IsEmpty() }
func (id CompID) IsEmpty() bool { return identity.ADT(id).IsEmpty() }

func (id DescID) String() string { return identity.ADT(id).String() }
func (id TypeID) String() string { return identity.ADT(id).String() }
func (id TermID) String() string { return identity.ADT(id).String() }
func (id CompID) String() string { return identity.ADT(id).String() }

func (rn DescRN) Next() DescRN { return DescRN(seqnum.ADT(rn).Next()) }
func (rn CompRN) Next() CompRN { return CompRN(seqnum.ADT(rn).Next()) }

// DescRef is the OCC boundary for every mutation of a Description and its
// owned Provider Type and Terms.
type DescRef struct {
	DescID DescID
	DescRN DescRN
}

// TypeRef and TermRef identify entities owned by a Description. They do not
// carry independent revisions; callers must use the owning DescRef for OCC.
type TypeRef struct {
	TypeID TypeID
}

type TermRef struct {
	TermID TermID
}

// CompRef is the independent OCC boundary for a Computation aggregate.
type CompRef struct {
	CompID CompID
	CompRN CompRN
}

// TypeResolver and TermResolver keep QN resolution statically typed. QNs share
// one namespace within proc, but Type and Term identities are not interchangeable.
type TypeResolver interface {
	ResolveTypeQN(context.Context, uniqsym.ADT) (TypeRef, error)
}

type TermResolver interface {
	ResolveTermQN(context.Context, uniqsym.ADT) (TermRef, error)
}
