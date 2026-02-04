package xactexp

import "orglang/go-engine/adt/uniqsym"

type ExpSpec interface {
	spec()
}

type OneSpec struct{}

func (OneSpec) spec() {}

type LinkSpec struct {
	XactQN uniqsym.ADT
}

// aka Internal Choice
type PlusSpec struct {
	Choices map[uniqsym.ADT]ExpSpec // conts
}

func (PlusSpec) spec() {}

// aka External Choice
type WithSpec struct {
	Choices map[uniqsym.ADT]ExpSpec // conts
}

func (WithSpec) spec() {}
