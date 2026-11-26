package adapter

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	typedef "orglang/go-engine/internal/pool/typedef/core"
)

func (dao *daoEnt) addExpRec(ctx context.Context, conn tx.Conn, rec typedef.ExpRec) error {
	ckAttr := slog.Any("ck", rec.Key())
	dao.log.Log(ctx, lf.LevelTrace, "insertion started", ckAttr)
	dto := dataFromExpRec(rec)
	txCtx := entpgx.WithConn(ctx, conn.Pgx)
	for _, st := range dto.States {
		_, err := dao.client.TypeExp.Create().
			SetID(st.ExpCK).
			SetSupExpCk(st.SupExpCK).
			SetKind(int16(st.K)).
			SetSpec(st.Spec).
			Save(txCtx)
		if err != nil {
			dao.log.Error("query execution failed", ckAttr, slog.Int64("exp_ck", st.ExpCK))
			return err
		}
	}
	dao.log.Log(ctx, lf.LevelTrace, "insertion succeed", ckAttr)
	return nil
}

func (dao *daoEnt) getRecByCK(ctx context.Context, conn tx.Conn, expCK contkey.ADT) (typedef.ExpRec, error) {
	ckAttr := slog.Any("ck", expCK)
	sql := selectRecByCK
	rows, err := conn.Pgx.Query(ctx, sql, contkey.ConvertToInt(expCK))
	if err != nil {
		dao.log.Error("query execution failed", ckAttr, slog.String("sql", sql))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[State])
	if err != nil {
		dao.log.Error("rows scanning failed", ckAttr)
		return nil, err
	}
	if len(dtos) == 0 { // revive:disable-line
		dao.log.Error("selection failed", ckAttr)
		return nil, errors.New("no rows selected")
	}
	dao.log.Log(ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	states := make(map[int64]State, len(dtos))
	for _, dto := range dtos {
		states[dto.ExpCK] = dto
	}
	return statesToExpRec(states, states[contkey.ConvertToInt(expCK)])
}

func (dao *daoEnt) getRecsByCKs(ctx context.Context, conn tx.Conn, expCKs []contkey.ADT) (_ []typedef.ExpRec, err error) {
	batch := pgx.Batch{}
	sql := selectRecByCK
	for _, expCK := range expCKs {
		batch.Queue(sql, contkey.ConvertToInt(expCK))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() { err = errors.Join(err, br.Close()) }()
	recs := make([]typedef.ExpRec, 0, len(expCKs))
	for _, expCK := range expCKs {
		ckAttr := slog.Any("ck", expCK)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", ckAttr, slog.String("sql", sql))
			return nil, readErr
		}
		dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[State])
		if scanErr != nil {
			dao.log.Error("rows scanning failed", ckAttr)
			return nil, scanErr
		}
		if len(dtos) == 0 {
			dao.log.Error("selection failed", ckAttr)
			return nil, typedef.ErrDoesNotExist(expCK)
		}
		rec, convErr := dataToExpRec(ExpRec{ExpCK: contkey.ConvertToInt(expCK), States: dtos})
		if convErr != nil {
			dao.log.Error("model conversion failed", ckAttr)
			return nil, convErr
		}
		recs = append(recs, rec)
	}
	dao.log.Log(ctx, lf.LevelTrace, "selection succeed", slog.Any("recs", recs))
	return recs, err
}

func (dao *daoEnt) GetTypeMap(ctx context.Context, conn tx.Conn, expCKs map[symbol.ADT]contkey.ADT) (_ map[symbol.ADT]typedef.ExpRec, err error) {
	batch := pgx.Batch{}
	sql := selectRecByCK
	for _, expCK := range expCKs {
		batch.Queue(sql, contkey.ConvertToInt(expCK))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() { err = errors.Join(err, br.Close()) }()
	recs := make(map[symbol.ADT]typedef.ExpRec, len(expCKs))
	for expPH, expCK := range expCKs {
		ckAttr := slog.Any("ck", expCK)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", ckAttr, slog.String("sql", sql))
			return nil, readErr
		}
		dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[State])
		if scanErr != nil {
			dao.log.Error("rows scanning failed", ckAttr)
			return nil, scanErr
		}
		dao.log.Log(ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
		if len(dtos) == 0 {
			dao.log.Error("selection failed", ckAttr)
			return nil, typedef.ErrDoesNotExist(expCK)
		}
		rec, convErr := dataToExpRec(ExpRec{ExpCK: contkey.ConvertToInt(expCK), States: dtos})
		if convErr != nil {
			dao.log.Error("model conversion failed", ckAttr)
			return nil, convErr
		}
		recs[expPH] = rec
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("recs", recs))
	return recs, err
}

const selectRecByCK = `
WITH RECURSIVE exp_tree AS (
    SELECT top.*
    FROM pool_type_exps top
    WHERE top.exp_ck = $1
    UNION ALL
    SELECT sub.*
    FROM pool_type_exps sub, exp_tree sup
    WHERE sub.sup_exp_ck = sup.exp_ck
)
SELECT * FROM exp_tree
`
