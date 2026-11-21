package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
)

type daoPgx struct {
	tableName string
	log       *slog.Logger
}

func NewDaoPgx(table string) func(log *slog.Logger) *daoPgx {
	return func(log *slog.Logger) *daoPgx {
		name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
		return &daoPgx{tableName: strings.TrimSpace(table), log: log.With(name)}
	}
}

// for compilation purposes
func newRepo() compsem.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) TouchRef(ctx context.Context, conn tx.Conn, ref compsem.SemRef) error {
	dto := DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	const query = `UPDATE %s SET comp_rn = comp_rn + 1 WHERE comp_id = $1 AND comp_rn = $2`
	sql := fmt.Sprintf(query, dao.tableName)
	ct, execErr := conn.Pgx.Exec(ctx, sql, dto.CompID, dto.CompRN)
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
