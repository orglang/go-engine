package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"orglang/go-engine/internal/adt/implsem"
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
func newRepo() implsem.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec implsem.SemRec) error {
	recAttr := slog.Any("rec", rec)
	dto := dataFromRec(rec)
	const query = `INSERT INTO %s (impl_qn, impl_id, kind) VALUES ($1, $2, $3)`
	sql := fmt.Sprintf(query, dao.tableName)
	_, execErr := conn.Pgx.Exec(ctx, sql, dto.ImplQN, dto.ImplID, dto.Kind)
	if execErr != nil {
		dao.log.Error("query execution failed", recAttr, slog.String("sql", sql))
		return execErr
	}
	return nil
}
