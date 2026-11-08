package adapter

import (
	"context"
	"log/slog"
	"reflect"

	descsem "orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/lib/tx"
)

type daoPgx struct {
	qb  queryBuilder
	log *slog.Logger
}

func NewDaoPgx(table string) func(log *slog.Logger) *daoPgx {
	return func(log *slog.Logger) *daoPgx {
		name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
		return &daoPgx{newSQLBuilder(table), log.With(name)}
	}
}

// for compilation purposes
func newRepo() descsem.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec descsem.SemRec) error {
	recAttr := slog.Any("rec", rec)
	dto, convErr := descsem.DataFromRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", recAttr)
		return convErr
	}
	sql, args := dao.qb.insertRec(dto)
	_, execErr := conn.Pgx.Exec(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", recAttr, slog.String("sql", sql))
		return execErr
	}
	return nil
}
