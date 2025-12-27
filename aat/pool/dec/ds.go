package dec

import (
	"orglang/orglang/avt/data"
)

// Port
type repo interface {
	Insert(data.Source, poolRec) error
}
