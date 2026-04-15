package compexec

import (
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/semcomp"
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

func (dao *pgxDAO) ModifyRec(db.Source, ExecMod) error {
	panic("unimplemented")
}

func (dao *pgxDAO) GetRecByRef(source db.Source, ref semcomp.CompRef) (ExecRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", ref)
	sql, args := dao.qb.selectRecByRef(semcomp.DataFromRef(ref))
	rows, execErr := ds.Conn.Query(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql), slog.Any("args", args))
		return ExecRec{}, execErr
	}
	defer rows.Close()
	dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[implRecDS])
	if scanErr != nil {
		dao.log.Error("rows scanning failed", refAttr, slog.Any("dto", dto))
		return ExecRec{}, scanErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dto", dto))
	rec, convErr := DataToImplRec(dto)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return ExecRec{}, convErr
	}
	return rec, nil
}
