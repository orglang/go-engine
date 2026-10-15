package adapter

import (
	"log/slog"
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/lib/db"
	"reflect"

	commexch "orglang/go-engine/proc/commexch/core"
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

func (dao *daoPgx) AddRec(db.UoW, commexch.ExchRec) error {
	panic("unimplemented")
}

func (dao *daoPgx) GetRefsByQNs(db.UoW, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error) {
	panic("unimplemented")
}

func (dao *daoPgx) GetSnapByQry(db.UoW, commexch.ExchQry) (commexch.ExchSnap, error) {
	panic("unimplemented")
}

func (dao *daoPgx) Modifyec(db.UoW, commexch.ExchMod) error {
	panic("unimplemented")
}
