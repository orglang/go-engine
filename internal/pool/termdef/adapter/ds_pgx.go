package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/uniqsym"

	termdef "orglang/go-engine/internal/pool/termdef/core"
)

type daoPgx struct {
	qb  termdef.QueryBuilder
	log *slog.Logger
}

func NewDaoPgx(qb termdef.QueryBuilder, log *slog.Logger) termdef.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() termdef.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec termdef.DefRec) error {
	refAttr := slog.Any("ref", rec.TermRef)
	dto, convErr := DataFromDecRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return convErr
	}
	sql, args := dao.qb.InsertRec(dto)
	_, execErr := conn.Pgx.Exec(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr)
		return execErr
	}
	return nil
}

func (dao *daoPgx) GetRecByQN(ctx context.Context, conn tx.Conn, qn uniqsym.ADT) (termdef.DefRec, error) {
	qnAttr := slog.Any("qn", qn)
	sql, args := dao.qb.SelectRecByQN(uniqsym.ConvertToString(qn))
	rows, execErr := conn.Pgx.Query(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", qnAttr)
		return termdef.DefRec{}, execErr
	}
	dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[termdef.DefRecDS])
	if scanErr != nil {
		dao.log.Error("rows scanning failed", qnAttr)
		return termdef.DefRec{}, scanErr
	}
	rec, convErr := DataToDecRec(dto)
	if convErr != nil {
		dao.log.Error("model conversion failed", qnAttr)
		return termdef.DefRec{}, convErr
	}
	return rec, nil
}
