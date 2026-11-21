package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"orglang/go-engine/internal/adt/descsem"
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
func newRepo() descsem.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec descsem.SemRec) error {
	recAttr := slog.Any("rec", rec)
	dto := dataFromRec(rec)
	const query = `INSERT INTO %s (desc_qn, desc_id, kind) VALUES ($1, $2, $3)`
	sql := fmt.Sprintf(query, dao.tableName)
	_, execErr := conn.Pgx.Exec(ctx, sql, dto.DescQN, dto.DescID, dto.Kind)
	if execErr != nil {
		dao.log.Error("query execution failed", recAttr, slog.String("sql", sql))
		return execErr
	}
	return nil
}
