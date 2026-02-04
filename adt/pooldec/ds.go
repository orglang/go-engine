package pooldec

import (
	"orglang/go-engine/lib/db"
)

// Port
type repo interface {
	Insert(db.Source, DecRec) error
}
