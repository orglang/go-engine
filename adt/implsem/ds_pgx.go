package implsem

import (
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"
)

type pgxDAO struct {
	qb  queryBuilder
	log *slog.Logger
}

func NewPgxDAO(table string) func(log *slog.Logger) *pgxDAO {
	return func(log *slog.Logger) *pgxDAO {
		name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
		return &pgxDAO{newSQLBuilder(table), log.With(name)}
	}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) AddRec(source db.Source, rec SemRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	recAttr := slog.Any("rec", rec)
	dto, convErr := DataFromRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", recAttr)
		return convErr
	}
	sql, args := dao.qb.insertRec(dto)
	_, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", recAttr, slog.String("sql", sql))
		return execErr
	}
	return nil
}
