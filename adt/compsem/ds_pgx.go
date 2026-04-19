package compsem

import (
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"
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

func (dao *pgxDAO) TouchRef(source db.Source, ref SemRef) error {
	ds := db.MustConform[db.SourcePgx](source)
	dto := DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	sql, args := dao.qb.updateRef(dto)
	ct, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	if ct.RowsAffected() == 0 {
		dao.log.Error("touching failed", refAttr)
		return errConcurrentModification(ref)
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "touching succeed", refAttr)
	return nil
}

func errConcurrentModification(got SemRef) error {
	return fmt.Errorf("concurrent modification: %v", got)
}
