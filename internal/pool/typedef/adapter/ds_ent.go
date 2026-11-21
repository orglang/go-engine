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
	refAttr := slog.Any("ref", rec.TypeRef)
	dao.log.Log(ctx, lf.LevelTrace, "addition started", refAttr)

	dto, err := DataFromDefRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}

	_, err = dao.client.TypeDef.Create().
		SetID(dto.TypeID).
		SetTypeRn(dto.TypeRN).
		SetExpCk(dto.ExpCK).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}

	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", refAttr)
	return nil
}

func (dao *daoEnt) GetRecsByQNs(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (_ map[uniqsym.ADT]typedef.DefRec, err error) {
	if len(typeQNs) == 0 {
		return map[uniqsym.ADT]typedef.DefRec{}, nil
	}
	const query = `SELECT td.type_id, td.type_rn, td.exp_ck
FROM pool_type_defs td
JOIN pool_desc_binds db ON db.desc_id = td.type_id
WHERE db.desc_qn = $1`
	batch := pgx.Batch{}
	for _, qn := range typeQNs {
		batch.Queue(query, uniqsym.ConvertToString(qn))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() { err = errors.Join(err, br.Close()) }()
	dtos := make(map[uniqsym.ADT]DefRec, len(typeQNs))
	for _, qn := range typeQNs {
		rows, readErr := br.Query()
		if readErr != nil {
			return nil, readErr
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
		if scanErr != nil {
			return nil, scanErr
		}
		dtos[qn] = dto
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return DataToDefRecMap(dtos)
}
