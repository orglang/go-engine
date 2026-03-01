package poolcomm

import (
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/poolcfg"
	"orglang/go-engine/adt/poolctx"
	"orglang/go-engine/adt/poolenv"
)

type pgxDAO struct {
	log *slog.Logger
}

func newPgxDAO(l *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{l.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) InsertRec(db.Source, CommRec) error {
	panic("unimplemented")
}

func (dao *pgxDAO) InsertRecs(db.Source, []CommRec) error {
	panic("unimplemented")
}

func (dao *pgxDAO) SelectCtxSnapByCtxSpec(db.Source, poolctx.CtxSpec) (poolctx.CtxSnap, error) {
	panic("unimplemented")
}

func (dao *pgxDAO) SelectEnvSnapByEnvSpec(db.Source, poolenv.EnvSpec) (poolenv.EnvSnap, error) {
	panic("unimplemented")
}

func (dao *pgxDAO) SelectCfgSnapBySpec(db.Source, poolcfg.CfgSpec) (poolcfg.CfgSnap, error) {
	panic("unimplemented")
}

func (dao *pgxDAO) SelectRecByRef(db.Source, commsem.SemRef) (CommRec, error) {
	panic("unimplemented")
}
