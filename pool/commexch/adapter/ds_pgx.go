package adapter

import (
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/uniqsym"

	commexch "orglang/go-engine/pool/commexch/core"

	"github.com/jackc/pgx/v5"
)

type daoPgx struct {
	qb  commexch.QueryBuilder
	log *slog.Logger
}

func NewDaoPgx(qb commexch.QueryBuilder, log *slog.Logger) commexch.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() commexch.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(ouw db.UoW, rec commexch.ExchRec) error {
	dto := DataFromRec(rec)
	refAttr := slog.Any("ref", rec.CommRef)
	sql, args := dao.qb.InsertRec(dto)
	_, execErr := ouw.Pgx.Exec(ouw.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(ouw.Ctx, lf.LevelTrace, "addition succeed", slog.Any("dto", dto))
	return nil
}

func (dao *daoPgx) ModifyRec(uow db.UoW, mod commexch.ExchMod) error {
	if mod.OffsetNr == nil {
		return nil
	}
	dto := DataFromMod(mod)
	refAttr := slog.Any("ref", mod.CommRef)
	sql, args := dao.qb.UpdateRec(dto)
	_, execErr := uow.Pgx.Exec(uow.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(uow.Ctx, lf.LevelTrace, "modification succeed", slog.Any("dto", dto))
	return nil
}

func (dao *daoPgx) GetRefsByQNs(uow db.UoW, qns []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error) {
	panic("unimplemented")
}

func (dao *daoPgx) GetSnapByQry(ds db.UoW, qry commexch.ExchQry) (commexch.ExchSnap, error) {
	refAttr := slog.Any("ref", qry.CommRef)
	qryDTO := DataFromQry(qry)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting started", slog.Any("qry", qryDTO))
	sql, args := dao.qb.SelectSnap(qryDTO)
	rows, execErr := ds.Pgx.Query(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return commexch.ExchSnap{}, execErr
	}
	defer rows.Close()
	snapDTO, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[commexch.ExchSnapDS])
	if scanErr != nil {
		dao.log.Error("rows scanning failed", refAttr)
		return commexch.ExchSnap{}, scanErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting succeed", slog.Any("dto", snapDTO))
	snap, convErr := DataToSnap(snapDTO)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return commexch.ExchSnap{}, convErr
	}
	return snap, nil
}
