package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/lib/tx"
)

type daoPgx struct {
	tableName string
	log       *slog.Logger
}

// NewDaoPgx creates a Description root repository for the supplied table.
func NewDaoPgx(table string) func(*slog.Logger) *daoPgx {
	return func(log *slog.Logger) *daoPgx {
		name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
		return &daoPgx{tableName: strings.TrimSpace(table), log: log.With(name)}
	}
}

func (dao *daoPgx) AddRef(ctx context.Context, conn tx.Conn, ref descroot.SemRef) error {
	const query = `
INSERT INTO %s (desc_id, desc_rn)
VALUES ($1, $2)
`
	sql := fmt.Sprintf(query, dao.tableName)
	_, err := conn.Pgx.Exec(ctx, sql, ref.DescID.String(), ref.DescRN)
	if err != nil {
		dao.log.Error("query execution failed", slog.Any("ref", ref), slog.String("sql", sql))
		return err
	}
	return nil
}

// TouchRef advances the Description revision if the expected revision matches.
func (dao *daoPgx) TouchRef(ctx context.Context, conn tx.Conn, ref descroot.SemRef) error {
	const query = `
UPDATE %s
SET desc_rn = desc_rn + 1
WHERE desc_id = $1 AND desc_rn = $2
`
	sql := fmt.Sprintf(query, dao.tableName)
	tag, err := conn.Pgx.Exec(ctx, sql, ref.DescID.String(), ref.DescRN)
	if err != nil {
		dao.log.Error("query execution failed", slog.Any("ref", ref), slog.String("sql", sql))
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("concurrent modification: %v", ref)
	}
	return nil
}
