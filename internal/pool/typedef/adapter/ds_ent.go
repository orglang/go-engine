package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	typedef "orglang/go-engine/internal/pool/typedef/core"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) typedef.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec typedef.DefRec) error {
	if err := dao.addExpRec(ctx, conn, rec.TypeExp); err != nil {
		return err
	}
	refAttr := slog.Any("typeID", rec.TypeID)
	dao.log.Log(ctx, lf.LevelTrace, "addition started", refAttr)

	dto, err := DataFromDefRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}

	_, err = dao.client.TypeDef.Create().
		SetID(dto.TypeID).
		SetProgID(dto.ProgID).
		SetExpCk(dto.ExpCK).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}

	const bindQuery = `
INSERT INTO pool_prog_binds (prog_qn, prog_id, kind)
VALUES ($1, $2, 1)
`
	if _, err := conn.Pgx.Exec(ctx, bindQuery, uniqsym.ConvertToString(rec.TypeQN), rec.TypeID.String()); err != nil {
		dao.log.Error("program binding insertion failed", refAttr, slog.Any("error", err))
		return err
	}

	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", refAttr)
	return nil
}

func (dao *daoEnt) GetRecsByQNs(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (_ map[uniqsym.ADT]typedef.DefRec, err error) {
	if len(typeQNs) == 0 {
		return map[uniqsym.ADT]typedef.DefRec{}, nil
	}
	const query = `
SELECT db.prog_qn AS type_qn, td.type_id, td.prog_id, dr.prog_rn AS prog_rn, td.exp_ck
FROM pool_type_defs AS td
JOIN pool_prog_roots AS dr
    ON dr.prog_id = td.prog_id
JOIN pool_prog_binds AS db
    ON db.prog_id = td.type_id
WHERE db.prog_qn = $1
`
	batch := pgx.Batch{}
	for _, qn := range typeQNs {
		batch.Queue(query, uniqsym.ConvertToString(qn))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	dtos := make(map[uniqsym.ADT]DefRec, len(typeQNs))
	for _, qn := range typeQNs {
		rows, readErr := br.Query()
		if readErr != nil {
			return nil, errors.Join(readErr, br.Close())
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
		if scanErr != nil {
			return nil, errors.Join(scanErr, br.Close())
		}
		dtos[qn] = dto
	}
	if closeErr := br.Close(); closeErr != nil {
		return nil, closeErr
	}
	recs, err := DataToDefRecMap(dtos)
	if err != nil {
		return nil, err
	}
	for qn, rec := range recs {
		rec.TypeQN = qn
		exp, err := dao.getRecByCK(ctx, conn, rec.ExpCK)
		if err != nil {
			return nil, err
		}
		rec.TypeExp = exp
		recs[qn] = rec
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("recs", recs))
	return recs, nil
}
