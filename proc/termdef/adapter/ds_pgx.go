package adapter

import (
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	termdef "orglang/go-engine/proc/termdef/core"
)

// Adapter
type daoPgx struct {
	log *slog.Logger
}

// for compilation purposes
func newRepo() termdef.Repo {
	return new(daoPgx)
}

func NewDaoPgx(l *slog.Logger) termdef.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{l.With(name)}
}

func (dao *daoPgx) InsertProc(db.UoW, termdef.DefRec) error {
	return nil
}
