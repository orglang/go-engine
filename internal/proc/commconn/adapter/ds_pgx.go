package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/commsem"
	"orglang/go-engine/internal/adt/uniqsym"

	commconn "orglang/go-engine/internal/proc/commconn/core"
)

type daoPgx struct {
	log *slog.Logger
}

func NewDaoPgx(log *slog.Logger) commconn.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{log.With(name)}
}

// for compilation purposes
func newRepo() commconn.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(context.Context, tx.Conn, commconn.ConnRec) error {
	panic("unimplemented")
}

func (dao *daoPgx) GetRefsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error) {
	panic("unimplemented")
}

func (dao *daoPgx) GetSnapByQry(context.Context, tx.Conn, commconn.ConnQry) (commconn.ConnSnap, error) {
	panic("unimplemented")
}

func (dao *daoPgx) Modifyec(context.Context, tx.Conn, commconn.ConnMod) error {
	panic("unimplemented")
}
