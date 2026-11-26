package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	typeexpent "orglang/go-engine/internal/proc/ent/typeexp"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

// Adapter
type expDaoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewExpDaoEnt(client *ent.Client, log *slog.Logger) typedef.TypeExpRepo {
	name := slog.String("name", reflect.TypeFor[expDaoEnt]().Name())
	return &expDaoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newExpRepo() typedef.TypeExpRepo {
	return new(expDaoEnt)
}

func (dao *expDaoEnt) AddRec(ctx context.Context, conn tx.Conn, rec typedef.ExpRec) error {
	ckAttr := slog.Any("expCK", rec.Key())
	dto := dataFromExpRec(rec)
	builders := make([]*ent.TypeExpCreate, 0, len(dto.States))
	for _, st := range dto.States {
		builders = append(builders, dao.client.TypeExp.Create().
			SetID(st.ExpCK).
			SetSupExpCk(st.SupExpCK).
			SetKind(int16(st.K)).
			SetSpec(st.Spec))
	}
	txCtx := entpgx.WithConn(ctx, conn.Pgx)
	err := dao.client.TypeExp.CreateBulk(builders...).
		OnConflictColumns(typeexpent.FieldID).
		DoNothing().
		Exec(txCtx)
	if err != nil {
		dao.log.Error("query execution failed", ckAttr)
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "insertion succeed", ckAttr)
	return nil
}

func (dao *expDaoEnt) SelectRecByCK(ctx context.Context, conn tx.Conn, expCK contkey.ADT) (typedef.ExpRec, error) {
	ckAttr := slog.Any("expCK", expCK)
	rows, queryErr := conn.Pgx.Query(ctx, selectByID, contkey.ConvertToInt(expCK))
	if queryErr != nil {
		dao.log.Error("query execution failed", ckAttr)
		return nil, queryErr
	}
	defer rows.Close()
	dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[State])
	if scanErr != nil {
		dao.log.Error("row scanning failed", ckAttr)
		return nil, scanErr
	}
	if len(dtos) == 0 {
		dao.log.Error("entity selection failed", ckAttr)
		return nil, typedef.ErrMissingInEnv(expCK)
	}
	dao.log.Log(ctx, lf.LevelTrace, "entity selection succeed", slog.Any("dtos", dtos))
	states := make(map[int64]State, len(dtos))
	for _, dto := range dtos {
		states[dto.ExpCK] = dto
	}
	return statesToExpRec(states, states[contkey.ConvertToInt(expCK)])
}

func (dao *expDaoEnt) SelectEnv(ctx context.Context, conn tx.Conn, expCKs []contkey.ADT) (map[contkey.ADT]typedef.ExpRec, error) {
	recs, err := dao.SelectRecsByCKs(ctx, conn, expCKs)
	if err != nil {
		return nil, err
	}
	env := make(map[contkey.ADT]typedef.ExpRec, len(recs))
	for _, rec := range recs {
		env[rec.Key()] = rec
	}
	return env, nil
}

func (dao *expDaoEnt) SelectRecsByCKs(ctx context.Context, conn tx.Conn, expCKs []contkey.ADT) (_ []typedef.ExpRec, err error) {
	batch := pgx.Batch{}
	for _, expCK := range expCKs {
		batch.Queue(selectByID, contkey.ConvertToInt(expCK))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	recs := make([]typedef.ExpRec, 0, len(expCKs))
	for _, expCK := range expCKs {
		ckAttr := slog.Any("expCK", expCK)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", ckAttr)
			return nil, readErr
		}
		dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[State])
		if scanErr != nil {
			dao.log.Error("rows scanning failed", ckAttr)
			return nil, scanErr
		}
		if len(dtos) == 0 {
			dao.log.Error("entity selection failed", ckAttr)
			return nil, typedef.ErrMissingInEnv(expCK)
		}
		rec, convErr := dataToExpRec(ExpRec{ExpCK: contkey.ConvertToInt(expCK), States: dtos})
		if convErr != nil {
			dao.log.Error("model conversion failed", ckAttr)
			return nil, convErr
		}
		recs = append(recs, rec)
	}
	dao.log.Log(ctx, lf.LevelTrace, "entities selection succeed", slog.Any("recs", recs))
	return recs, err
}

const selectByID = `
	WITH RECURSIVE exp_tree AS (
		SELECT top.*
		FROM proc_type_exps top
		WHERE top.exp_ck = $1
		UNION ALL
		SELECT sub.*
		FROM proc_type_exps sub, exp_tree sup
		WHERE sub.sup_exp_ck = sup.exp_ck
	)
	SELECT * FROM exp_tree`
