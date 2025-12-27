package alias

import (
	"orglang/orglang/avt/data"
)

type Repo interface {
	Insert(data.Source, Root) error
}

type rootDS struct {
	ID  string
	RN  int64
	Sym string
}
