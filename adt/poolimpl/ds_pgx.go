package poolimpl

import (
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/implsem"
)

type pgxDAO struct {
	qb  queryBuilder
	log *slog.Logger
}

func newPgxDAO(qb queryBuilder, log *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{qb, log.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) GetRecByRef(source db.Source, ref implsem.SemRef) (ImplRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", ref)
	sql, args := dao.qb.selectRecByRef(implsem.DataFromRef(ref))
	rows, execErr := ds.Conn.Query(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql), slog.Any("args", args))
		return ImplRec{}, execErr
	}
	defer rows.Close()
	dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[implRecDS])
	if scanErr != nil {
		dao.log.Error("rows scanning failed", refAttr, slog.Any("dto", dto))
		return ImplRec{}, scanErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dto", dto))
	rec, convErr := DataToImplRec(dto)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return ImplRec{}, convErr
	}
	return rec, nil
}
