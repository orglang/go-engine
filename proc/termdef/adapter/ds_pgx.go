package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/tx"

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

func (dao *daoPgx) InsertProc(context.Context, tx.Conn, termdef.DefRec) error {
	return nil
}
