package identity

import (
	"errors"

	"github.com/rs/xid"
)

var (
	Nil ADT
)

type Identifiable interface {
	Ident() ADT
}

type ADT xid.ID

func (ADT) PH() {}

func New() ADT {
	return ADT(xid.New())
}

func Empty() ADT {
	return ADT(xid.NilID())
}

func (id ADT) IsEmpty() bool {
	return xid.ID(id).IsZero()
}

func (id ADT) String() string {
	return xid.ID(id).String()
}

var (
	ErrEmpty error = errors.New("empty id")
)
