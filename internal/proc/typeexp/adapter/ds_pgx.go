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

	typeexp "orglang/go-engine/internal/proc/typeexp/core"
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

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec typeexp.ExpRec) (err error) {
	ckAttr := slog.Any("expCK", rec.Key())
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

func (dao *daoPgx) SelectRecByCK(ctx context.Context, conn tx.Conn, expCK contkey.ADT) (typeexp.ExpRec, error) {
	ckAttr := slog.Any("expCK", expCK)
	rows, queryErr := conn.Pgx.Query(ctx, selectByID, contkey.ConvertToInt(expCK))
	if queryErr != nil {
		dao.log.Error("query execution failed", ckAttr)
		return nil, queryErr
	}
	defer rows.Close()
	dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[typeexp.StateDS])
	if scanErr != nil {
		dao.log.Error("row scanning failed", ckAttr)
		return nil, scanErr
	}
	if len(dtos) == 0 { // revive:disable-line
		dao.log.Error("entity selection failed", ckAttr)
		return nil, errors.New("no rows selected")
	}
	dao.log.Log(ctx, lf.LevelTrace, "entity selection succeed", slog.Any("dtos", dtos))
	states := make(map[int64]typeexp.StateDS, len(dtos))
	for _, dto := range dtos {
		states[dto.ExpCK] = dto
	}
	return statesToExpRec(states, states[contkey.ConvertToInt(expCK)])
}

func (dao *daoPgx) SelectEnv(ctx context.Context, conn tx.Conn, expCKs []contkey.ADT) (map[contkey.ADT]typeexp.ExpRec, error) {
	recs, err := dao.SelectRecsByCKs(ctx, conn, expCKs)
	if err != nil {
		return nil, err
	}
	env := make(map[contkey.ADT]typeexp.ExpRec, len(recs))
	for _, rec := range recs {
		env[rec.Key()] = rec
	}
	return env, nil
}

func (dao *daoPgx) SelectRecsByCKs(ctx context.Context, conn tx.Conn, expCKs []contkey.ADT) (_ []typeexp.ExpRec, err error) {
	batch := pgx.Batch{}
	for _, expCK := range expCKs {
		batch.Queue(selectByID, contkey.ConvertToInt(expCK))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var recs []typeexp.ExpRec
	for _, expCK := range expCKs {
		ckAttr := slog.Any("expCK", expCK)
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", ckAttr)
		}
		dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[typeexp.StateDS])
		if err != nil {
			dao.log.Error("rows scanning failed", ckAttr)
		}
		if len(dtos) == 0 {
			dao.log.Error("entity selection failed", ckAttr)
			return nil, typeexp.ErrDoesNotExist(expCK)
		}
		rec, err := dataToExpRec(typeexp.ExpRecDS{ExpCK: contkey.ConvertToInt(expCK), States: dtos})
		if err != nil {
			dao.log.Error("model conversion failed", ckAttr)
			return nil, err
		}
		recs = append(recs, rec)
	}
	dao.log.Log(ctx, lf.LevelTrace, "entities selection succeed", slog.Any("recs", recs))
	return recs, err
}

const (
	insertRec = `
		insert into proc_type_exps (
			exp_ck, sup_exp_ck, desc_id, desc_rn, kind, spec
		) values (
			@exp_ck, @sup_exp_ck, @desc_id, @desc_rn, @kind, @spec
		)
		on conflict (exp_ck) do nothing`

	selectByID = `
		with recursive exp_tree AS (
			select top.*
			from proc_type_exps top
			where exp_ck = $1
			union all
			select sub.*
			from proc_type_exps sub, exp_tree sup
			where sub.sup_exp_ck = sup.exp_ck
		)
		select * from exp_tree`
)
