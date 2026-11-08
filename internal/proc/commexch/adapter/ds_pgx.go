package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/commsem"
	"orglang/go-engine/internal/adt/uniqsym"

	commexch "orglang/go-engine/internal/proc/commexch/core"
)

type daoPgx struct {
	log *slog.Logger
}

func NewDaoPgx(log *slog.Logger) commexch.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{log.With(name)}
}

// for compilation purposes
func newRepo() commexch.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(context.Context, tx.Conn, commexch.ExchRec) error {
	panic("unimplemented")
}

func (dao *daoPgx) GetRefsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error) {
	panic("unimplemented")
}

func (dao *daoPgx) GetSnapByQry(context.Context, tx.Conn, commexch.ExchQry) (commexch.ExchSnap, error) {
	panic("unimplemented")
}

func (dao *daoPgx) Modifyec(context.Context, tx.Conn, commexch.ExchMod) error {
	panic("unimplemented")
}
