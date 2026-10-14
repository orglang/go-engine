package adapter

import (
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"

	termdef "orglang/go-engine/pool/termdef/core"
)

type daoPgx struct {
	qb  termdef.QueryBuilder
	log *slog.Logger
}

func newDaoPgx(qb termdef.QueryBuilder, log *slog.Logger) *daoPgx {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() termdef.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(uow db.UoW, rec termdef.DefRec) error {
	refAttr := slog.Any("ref", rec.TermRef)
	dto, convErr := DataFromDecRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return convErr
	}
	sql, args := dao.qb.InsertRec(dto)
	_, execErr := uow.Pgx.Exec(uow.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr)
		return execErr
	}
	return nil
}

func (dao *daoPgx) GetRecByQN(uow db.UoW, qn uniqsym.ADT) (termdef.DefRec, error) {
	qnAttr := slog.Any("qn", qn)
	sql, args := dao.qb.SelectRecByQN(uniqsym.ConvertToString(qn))
	rows, execErr := uow.Pgx.Query(uow.Ctx, sql, args...)
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
