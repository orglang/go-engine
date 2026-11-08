package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/lf"
	"orglang/go-engine/lib/tx"

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

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec commexch.ExchRec) error {
	dto := DataFromRec(rec)
	refAttr := slog.Any("ref", rec.CommRef)
	sql, args := dao.qb.InsertRec(dto)
	_, execErr := conn.Pgx.Exec(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", slog.Any("dto", dto))
	return nil
}

func (dao *daoPgx) ModifyRec(ctx context.Context, conn tx.Conn, mod commexch.ExchMod) error {
	if mod.OffsetNr == nil {
		return nil
	}
	dto := DataFromMod(mod)
	refAttr := slog.Any("ref", mod.CommRef)
	sql, args := dao.qb.UpdateRec(dto)
	_, execErr := conn.Pgx.Exec(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(ctx, lf.LevelTrace, "modification succeed", slog.Any("dto", dto))
	return nil
}

func (dao *daoPgx) GetRefsByQNs(ctx context.Context, conn tx.Conn, qns []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error) {
	panic("unimplemented")
}

func (dao *daoPgx) GetSnapByQry(ctx context.Context, conn tx.Conn, qry commexch.ExchQry) (commexch.ExchSnap, error) {
	refAttr := slog.Any("ref", qry.CommRef)
	qryDTO := DataFromQry(qry)
	dao.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("qry", qryDTO))
	sql, args := dao.qb.SelectSnap(qryDTO)
	rows, execErr := conn.Pgx.Query(ctx, sql, args...)
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
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dto", snapDTO))
	snap, convErr := DataToSnap(snapDTO)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return commexch.ExchSnap{}, convErr
	}
	return snap, nil
}
