package adapter

import (
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	proctermdef "orglang/go-engine/proc/termdef/core"
)

// Adapter
type daoPgx struct {
	log *slog.Logger
}

// for compilation purposes
func newRepo() proctermdef.Repo {
	return new(daoPgx)
}

func NewDaoPgx(l *slog.Logger) proctermdef.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{l.With(name)}
}

func (dao *daoPgx) InsertProc(db.UoW, proctermdef.DefRec) error {
	return nil
}
