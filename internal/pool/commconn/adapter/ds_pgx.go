package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	commconn "orglang/go-engine/internal/pool/commconn/core"
	commturnadapter "orglang/go-engine/internal/pool/commturn/adapter"
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
func newRepo() commconn.Repo { return new(daoPgx) }

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec commconn.ConnRec) error {
	dto := DataFromRec(rec)
	refAttr := slog.Any("conn_id", rec.ConnID)
	sql, args := dao.qb.InsertRec(dto)
	_, err := conn.Pgx.Exec(ctx, sql, args...)
	if err != nil {
		dao.log.Error("Connection insertion failed", refAttr, slog.String("sql", sql))
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection insertion succeeded", slog.Any("dto", dto))
	return nil
}

func (dao *daoPgx) ModifyRec(ctx context.Context, conn tx.Conn, mod commconn.ConnMod) error {
	refAttr := slog.Any("conn_id", mod.ConnID)
	if mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty() {
		dto := DataFromMod(mod)
		sql, args := dao.qb.UpdateRec(dto)
		tag, err := conn.Pgx.Exec(ctx, sql, args...)
		if err != nil {
			dao.log.Error("Connection offset update failed", refAttr, slog.String("sql", sql))
			return err
		}
		if tag.RowsAffected() == 0 {
			return commconn.ErrSnapshotConflict(mod.PrevOffsetNr)
		}
	}
	for _, turn := range mod.Turns {
		dto := commturnadapter.DataFromStepRec(turn)
		const nextAct = `
			UPDATE pool_comp_conns
			SET act_nr = act_nr + 1
			WHERE conn_id = $1
			RETURNING act_nr
		`
		if err := conn.Pgx.QueryRow(ctx, nextAct, identity.ConvertToString(turn.GetConnID())).Scan(&dto.ActNr); err != nil {
			dao.log.Error("Connection Act sequence update failed", refAttr)
			return err
		}
		sql, args := dao.qb.InsertTurnRec(dto)
		if _, err := conn.Pgx.Exec(ctx, sql, args...); err != nil {
			dao.log.Error("Act insertion failed", refAttr, slog.String("sql", sql))
			return err
		}
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection modification succeeded", refAttr)
	return nil
}

func (dao *daoPgx) GetSnapByQry(ctx context.Context, conn tx.Conn, qry commconn.ConnQry) (commconn.ConnSnap, error) {
	refAttr := slog.Any("conn_id", qry.ConnID)
	qryDTO := DataFromQry(qry)
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot selection started", slog.Any("qry", qryDTO))
	sql, args := dao.qb.SelectSnap(qryDTO)
	rows, err := conn.Pgx.Query(ctx, sql, args...)
	if err != nil {
		dao.log.Error("Connection snapshot query failed", refAttr, slog.String("sql", sql))
		return commconn.ConnSnap{}, err
	}
	defer rows.Close()
	snapDTO, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[commconn.ConnSnapDS])
	if err != nil {
		dao.log.Error("Connection snapshot scan failed", refAttr)
		return commconn.ConnSnap{}, err
	}
	snap, err := DataToSnap(snapDTO)
	if err != nil {
		dao.log.Error("Connection snapshot conversion failed", refAttr)
		return commconn.ConnSnap{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot loaded", slog.Any("dto", snapDTO))
	return snap, nil
}
