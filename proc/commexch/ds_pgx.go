package commexch

import (
	"log/slog"
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/lib/db"
	"reflect"
)

type pgxDAO struct {
	log *slog.Logger
}

func newPgxDAO(log *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{log.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) AddRec(db.Source, ExchRec) error {
	panic("unimplemented")
}

func (dao *pgxDAO) GetRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error) {
	panic("unimplemented")
}

func (dao *pgxDAO) GetSnapByQry(db.Source, ExchQry) (ExchSnap, error) {
	panic("unimplemented")
}

func (dao *pgxDAO) Modifyec(db.Source, ExchMod) error {
	panic("unimplemented")
}
