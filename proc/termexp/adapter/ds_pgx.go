package adapter

import (
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	proctermexp "orglang/go-engine/proc/termexp/core"
)

// Adapter
type daoPgx struct {
	log *slog.Logger
}

// for compilation purposes
func newRepo() proctermexp.Repo {
	return new(daoPgx)
}

func NewDaoPgx(l *slog.Logger) proctermexp.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{l.With(name)}
}

func (dao *daoPgx) Insert(uow db.UoW, rec proctermexp.ExpRec) error {
	return nil
}
