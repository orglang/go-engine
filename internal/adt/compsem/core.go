// Package compsem is a compatibility facade for the comproot aggregate reference.
package compsem

import "orglang/go-engine/internal/adt/comproot"

type SemRef = comproot.SemRef

func New() SemRef { return comproot.New() }
