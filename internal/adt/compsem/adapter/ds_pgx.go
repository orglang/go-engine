package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	compsem "orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/lib/lf"
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
func newRepo() compsem.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) TouchRef(ctx context.Context, conn tx.Conn, ref compsem.SemRef) error {
	dto := compsem.DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	sql, args := dao.qb.updateRef(dto)
	ct, execErr := conn.Pgx.Exec(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	if ct.RowsAffected() == 0 {
		dao.log.Error("touching failed", refAttr)
		return errConcurrentModification(ref)
	}
	dao.log.Log(ctx, lf.LevelTrace, "touching succeed", refAttr)
	return nil
}

func errConcurrentModification(got compsem.SemRef) error {
	return fmt.Errorf("concurrent modification: %v", got)
}
