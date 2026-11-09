package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/typesem"

	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

// Adapter
type daoPgx struct {
	qb  typeexp.QueryBuilder
	log *slog.Logger
}

func NewDaoPgx(qb typeexp.QueryBuilder, log *slog.Logger) typeexp.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() typeexp.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec typeexp.ExpRec, ref typesem.SemRef) (err error) {
	ckAttr := slog.Any("ck", rec.Key())
	dto := dataFromExpRec(rec)
	batch := pgx.Batch{}
	for _, st := range dto.States {
		sql, args := dao.qb.InsertRec(st)
		batch.Queue(sql, args...)
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for range dto.States {
		_, readErr := br.Exec()
		if readErr != nil {
			dao.log.Error("query execution failed", ckAttr)
			return readErr
		}
	}
	return nil
}

func (dao *daoPgx) GetRecByCK(ctx context.Context, conn tx.Conn, expCK contkey.ADT) (typeexp.ExpRec, error) {
	ckAttr := slog.Any("ck", expCK)
	sql := dao.qb.SelectRecByCK()
	rows, err := conn.Pgx.Query(ctx, sql, contkey.ConvertToInt(expCK))
	if err != nil {
		dao.log.Error("query execution failed", ckAttr, slog.String("sql", sql))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[typeexp.StateDS])
	if err != nil {
		dao.log.Error("rows scanning failed", ckAttr)
		return nil, err
	}
	if len(dtos) == 0 { // revive:disable-line
		dao.log.Error("selection failed", ckAttr)
		return nil, errors.New("no rows selected")
	}
	dao.log.Log(ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	states := make(map[int64]typeexp.StateDS, len(dtos))
	for _, dto := range dtos {
		states[dto.ExpCK] = dto
	}
	return statesToExpRec(states, states[contkey.ConvertToInt(expCK)])
}

func (dao *daoPgx) GetRecsByCKs(ctx context.Context, conn tx.Conn, expCKs []contkey.ADT) (_ []typeexp.ExpRec, err error) {
	batch := pgx.Batch{}
	sql := dao.qb.SelectRecByCK()
	for _, expCK := range expCKs {
		batch.Queue(sql, contkey.ConvertToInt(expCK))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	recs := make([]typeexp.ExpRec, 0, len(expCKs))
	for _, expCK := range expCKs {
		ckAttr := slog.Any("ck", expCK)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", ckAttr, slog.String("sql", sql))
			return nil, readErr
		}
		dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[typeexp.StateDS])
		if scanErr != nil {
			dao.log.Error("rows scanning failed", ckAttr)
			return nil, scanErr
		}
		if len(dtos) == 0 {
			dao.log.Error("selection failed", ckAttr)
			return nil, typeexp.ErrDoesNotExist(expCK)
		}
		rec, convErr := dataToExpRec(typeexp.ExpRecDS{ExpCK: contkey.ConvertToInt(expCK), States: dtos})
		if convErr != nil {
			dao.log.Error("model conversion failed", ckAttr)
			return nil, convErr
		}
		recs = append(recs, rec)
	}
	dao.log.Log(ctx, lf.LevelTrace, "selection succeed", slog.Any("recs", recs))
	return recs, err
}

func (dao *daoPgx) GetRecMap(ctx context.Context, conn tx.Conn, expCKs map[symbol.ADT]contkey.ADT) (_ map[symbol.ADT]typeexp.ExpRec, err error) {
	batch := pgx.Batch{}
	sql := dao.qb.SelectRecByCK()
	for _, expCK := range expCKs {
		batch.Queue(sql, contkey.ConvertToInt(expCK))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	recs := make(map[symbol.ADT]typeexp.ExpRec, len(expCKs))
	for expPH, expCK := range expCKs {
		ckAttr := slog.Any("ck", expCK)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", ckAttr, slog.String("sql", sql))
			return nil, readErr
		}
		dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[typeexp.StateDS])
		if scanErr != nil {
			dao.log.Error("rows scanning failed", ckAttr)
			return nil, scanErr
		}
		dao.log.Log(ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
		if len(dtos) == 0 {
			dao.log.Error("selection failed", ckAttr)
			return nil, typeexp.ErrDoesNotExist(expCK)
		}
		rec, convErr := dataToExpRec(typeexp.ExpRecDS{ExpCK: contkey.ConvertToInt(expCK), States: dtos})
		if convErr != nil {
			dao.log.Error("model conversion failed", ckAttr)
			return nil, convErr
		}
		recs[expPH] = rec
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("recs", recs))
	return recs, err
}
