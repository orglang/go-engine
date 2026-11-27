package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"

	"orglang/go-engine/internal/proc/ent"
	typedefent "orglang/go-engine/internal/proc/ent/typedef"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type daoEnt struct { client *ent.Client; log *slog.Logger }

func NewDaoEnt(client *ent.Client, log *slog.Logger) typedef.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}
func newRepo() typedef.Repo { return new(daoEnt) }

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec typedef.DefRec) error {
	idAttr := slog.Any("typeID", rec.TypeID)
	dao.log.Log(ctx, lf.LevelTrace, "addition started", idAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil { dao.log.Error("model conversion failed", idAttr); return err }
	_, err = dao.client.TypeDef.Create().SetID(dto.TypeID).SetDescID(dto.DescID).SetExpCk(dto.ExpCK).Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil { dao.log.Error("query execution failed", idAttr); return err }
	const bindQuery = `
INSERT INTO proc_desc_binds (desc_qn, desc_id, kind)
VALUES ($1, $2, 1)
`
	if _, err := conn.Pgx.Exec(ctx, bindQuery, uniqsym.ConvertToString(rec.TypeQN), rec.TypeID.String()); err != nil {
		dao.log.Error("description binding insertion failed", idAttr, slog.Any("error", err)); return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", idAttr)
	return nil
}

func (dao *daoEnt) ModifyRec(ctx context.Context, conn tx.Conn, rec typedef.DefRec) error {
	idAttr := slog.Any("typeID", rec.TypeID)
	dao.log.Log(ctx, lf.LevelTrace, "modification started", idAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil { dao.log.Error("model conversion failed", idAttr); return err }
	affected, err := dao.client.TypeDef.Update().Where(typedefent.ID(dto.TypeID)).SetExpCk(dto.ExpCK).Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil { dao.log.Error("query execution failed", idAttr); return err }
	if affected == 0 { dao.log.Error("entity update failed", idAttr); return typedef.ErrDoesNotExist(rec.TypeID) }
	dao.log.Log(ctx, lf.LevelTrace, "modification succeed", idAttr)
	return nil
}

func (dao *daoEnt) GetIDs(ctx context.Context, conn tx.Conn) ([]identity.ADT, error) {
	rows, err := conn.Pgx.Query(ctx, selectIDs)
	if err != nil { dao.log.Error("query execution failed"); return nil, err }
	defer rows.Close()
	values, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil { dao.log.Error("rows scanning failed"); return nil, err }
	ids := make([]identity.ADT, len(values))
	for i, value := range values {
		ids[i], err = identity.ConvertFromString(value)
		if err != nil { return nil, err }
	}
	return ids, nil
}

func (dao *daoEnt) GetRecByID(ctx context.Context, conn tx.Conn, id identity.ADT) (typedef.DefRec, error) {
	rows, err := conn.Pgx.Query(ctx, selectRecByID, id.String())
	if err != nil { dao.log.Error("query execution failed", slog.Any("typeID", id)); return typedef.DefRec{}, err }
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
	if err != nil { dao.log.Error("row scanning failed", slog.Any("typeID", id)); return typedef.DefRec{}, err }
	return DataToDefRec(dto)
}

func (dao *daoEnt) GetRecByQN(ctx context.Context, conn tx.Conn, qn uniqsym.ADT) (typedef.DefRec, error) {
	rows, err := conn.Pgx.Query(ctx, selectRecByQN, uniqsym.ConvertToString(qn))
	if err != nil { dao.log.Error("query execution failed", slog.Any("typeQN", qn)); return typedef.DefRec{}, err }
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
	if err != nil { dao.log.Error("row scanning failed", slog.Any("typeQN", qn)); return typedef.DefRec{}, err }
	return DataToDefRec(dto)
}

func (dao *daoEnt) GetRecsByIDs(ctx context.Context, conn tx.Conn, ids []identity.ADT) (_ []typedef.DefRec, err error) {
	if len(ids) == 0 { return []typedef.DefRec{}, nil }
	batch := pgx.Batch{}
	for _, id := range ids {
		if id.IsEmpty() { return nil, identity.ErrEmpty }
		batch.Queue(selectRecByID, id.String())
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() { err = errors.Join(err, br.Close()) }()
	dtos := make([]DefRec, 0, len(ids))
	for _, id := range ids {
		rows, queryErr := br.Query()
		if queryErr != nil { dao.log.Error("query execution failed", slog.Any("typeID", id)); return nil, queryErr }
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
		if scanErr != nil { dao.log.Error("row scanning failed", slog.Any("typeID", id)); return nil, scanErr }
		dtos = append(dtos, dto)
	}
	if err != nil { return nil, err }
	return DataToDefRecs(dtos)
}

func (dao *daoEnt) SelectEnv(ctx context.Context, conn tx.Conn, qns []uniqsym.ADT) (map[uniqsym.ADT]typedef.DefRec, error) {
	recs, err := dao.GetRecsByQNs(ctx, conn, qns)
	if err != nil { return nil, err }
	env := make(map[uniqsym.ADT]typedef.DefRec, len(recs))
	for i, rec := range recs { env[qns[i]] = rec }
	return env, nil
}

func (dao *daoEnt) GetRecsByQNs(ctx context.Context, conn tx.Conn, qns []uniqsym.ADT) (_ []typedef.DefRec, err error) {
	if len(qns) == 0 { return []typedef.DefRec{}, nil }
	batch := pgx.Batch{}
	for _, qn := range qns { batch.Queue(selectRecByQN, uniqsym.ConvertToString(qn)) }
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() { err = errors.Join(err, br.Close()) }()
	dtos := make([]DefRec, 0, len(qns))
	for _, qn := range qns {
		rows, queryErr := br.Query()
		if queryErr != nil { dao.log.Error("query execution failed", slog.Any("typeQN", qn)); return nil, queryErr }
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
		if scanErr != nil { dao.log.Error("row scanning failed", slog.Any("typeQN", qn)); return nil, scanErr }
		dtos = append(dtos, dto)
	}
	if err != nil { return nil, err }
	return DataToDefRecs(dtos)
}

const (
	selectRecByQN = `
SELECT db.desc_qn AS type_qn, td.type_id, td.desc_id, td.exp_ck, dr.desc_rn AS desc_rn
FROM proc_type_defs AS td
JOIN proc_desc_roots AS dr ON dr.desc_id = td.desc_id
JOIN proc_desc_binds AS db ON db.desc_id = td.type_id
WHERE db.desc_qn = $1
`
	selectRecByID = `
SELECT db.desc_qn AS type_qn, td.type_id, td.desc_id, td.exp_ck, dr.desc_rn AS desc_rn
FROM proc_type_defs AS td
JOIN proc_desc_roots AS dr ON dr.desc_id = td.desc_id
JOIN proc_desc_binds AS db ON db.desc_id = td.type_id
WHERE td.type_id = $1
`
	selectIDs = `
SELECT td.type_id
FROM proc_type_defs AS td
`
)
