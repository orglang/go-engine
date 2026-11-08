package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/seqnum"

	commconn "orglang/go-engine/internal/pool/commconn/core"
	commturnadapter "orglang/go-engine/internal/pool/commturn/adapter"

	"github.com/jackc/pgx/v5"
)

type daoPgx struct {
	qb  commconn.QueryBuilder
	log *slog.Logger
}

func NewDaoPgx(qb commconn.QueryBuilder, log *slog.Logger) commconn.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() commconn.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec commconn.ConnRec) error {
	dto := DataFromRec(rec)
	refAttr := slog.Any("ref", rec.ConnRef)
	sql, args := dao.qb.InsertRec(dto)
	_, execErr := conn.Pgx.Exec(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", slog.Any("dto", dto))
	return nil
}

func (dao *daoPgx) ModifyRec(ctx context.Context, conn tx.Conn, mod commconn.ConnMod) error {
	refAttr := slog.Any("ref", mod.ConnRef)
	if mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty() {
		dto := DataFromMod(mod)
		sql, args := dao.qb.UpdateRec(dto)
		tag, execErr := conn.Pgx.Exec(ctx, sql, args...)
		if execErr != nil {
			dao.log.Error("Connection offset update failed", refAttr, slog.String("sql", sql))
			return execErr
		}
		if tag.RowsAffected() == 0 {
			return commconn.ErrSnapshotConflict(mod.ConnRef.ConnRN, mod.PrevOffsetNr)
		}
	}

	nextRNs := make(map[string]seqnum.ADT)
	for _, turn := range mod.Turns {
		ref := turn.GetConnRef()
		key := ref.ConnID.String()
		expectedRN, ok := nextRNs[key]
		if !ok {
			expectedRN = ref.ConnRN
		}
		const updateRevision = `
			UPDATE pool_connections
			SET conn_rn = conn_rn + 1
			WHERE conn_id = $1 AND conn_rn = $2
		`
		tag, err := conn.Pgx.Exec(ctx, updateRevision, key, int64(expectedRN))
		if err != nil {
			dao.log.Error("Connection revision update failed", refAttr)
			return err
		}
		if tag.RowsAffected() == 0 {
			return commconn.ErrRevisionConflict(expectedRN)
		}
		dto := commturnadapter.DataFromStepRec(turn)
		dto.ConnRN = int64(expectedRN)
		sql, args := dao.qb.InsertTurnRec(dto)
		if _, execErr := conn.Pgx.Exec(ctx, sql, args...); execErr != nil {
			dao.log.Error("Act insertion failed", refAttr, slog.String("sql", sql))
			return execErr
		}
		nextRNs[key] = expectedRN.Next()
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection modification succeed", refAttr)
	return nil
}

func (dao *daoPgx) GetSnapByQry(ctx context.Context, conn tx.Conn, qry commconn.ConnQry) (commconn.ConnSnap, error) {
	refAttr := slog.Any("ref", qry.ConnRef)
	qryDTO := DataFromQry(qry)
	dao.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("qry", qryDTO))
	sql, args := dao.qb.SelectSnap(qryDTO)
	rows, execErr := conn.Pgx.Query(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return commconn.ConnSnap{}, execErr
	}
	defer rows.Close()
	snapDTO, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[commconn.ConnSnapDS])
	if scanErr != nil {
		dao.log.Error("rows scanning failed", refAttr)
		return commconn.ConnSnap{}, scanErr
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dto", snapDTO))
	snap, convErr := DataToSnap(snapDTO)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return commconn.ConnSnap{}, convErr
	}
	return snap, nil
}
