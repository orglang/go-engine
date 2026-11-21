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
	"orglang/go-engine/internal/adt/typesem"
	typesemadapter "orglang/go-engine/internal/adt/typesem/adapter"
	"orglang/go-engine/internal/adt/uniqsym"

	"orglang/go-engine/internal/proc/ent"
	typedefent "orglang/go-engine/internal/proc/ent/typedef"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) typedef.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() typedef.Repo {
	return new(daoEnt)
}

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec typedef.DefRec) error {
	idAttr := slog.Any("typeID", rec.TypeRef.TypeID)
	dao.log.Log(ctx, lf.LevelTrace, "addition started", idAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", idAttr)
		return err
	}
	_, err = dao.client.TypeDef.Create().
		SetID(dto.TypeID).
		SetTypeRn(dto.TypeRN).
		SetExpCk(dto.ExpCK).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", idAttr)
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", idAttr)
	return nil
}

func (dao *daoEnt) ModifyRec(ctx context.Context, conn tx.Conn, rec typedef.DefRec) error {
	idAttr := slog.Any("typeID", rec.TypeRef.TypeID)
	dao.log.Log(ctx, lf.LevelTrace, "modification started", idAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", idAttr)
		return err
	}
	affected, err := dao.client.TypeDef.Update().
		Where(typedefent.ID(dto.TypeID), typedefent.TypeRn(dto.TypeRN-1)).
		SetTypeRn(dto.TypeRN).
		SetExpCk(dto.ExpCK).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", idAttr)
		return err
	}
	if affected == 0 {
		dao.log.Error("entity update failed", idAttr)
		return typedef.ErrOptimisticUpdate(rec.TypeRef.TypeRN - 1)
	}
	dao.log.Log(ctx, lf.LevelTrace, "modification succeed", idAttr)
	return nil
}

func (dao *daoEnt) GetRefs(ctx context.Context, conn tx.Conn) ([]typesem.SemRef, error) {
	rows, err := conn.Pgx.Query(ctx, selectRefs)
	if err != nil {
		dao.log.Error("query execution failed")
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[typesemadapter.SemRef])
	if err != nil {
		dao.log.Error("rows scanning failed")
		return nil, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return typesemadapter.DataToRefs(dtos)
}

func (dao *daoEnt) GetRecByRef(ctx context.Context, conn tx.Conn, ref typesem.SemRef) (typedef.DefRec, error) {
	refAttr := slog.Any("ref", ref)
	rows, err := conn.Pgx.Query(ctx, selectRecByID, ref.TypeID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return typedef.DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
	if err != nil {
		dao.log.Error("row scanning failed", refAttr)
		return typedef.DefRec{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entity selection succeed", refAttr)
	return DataToDefRec(dto)
}

func (dao *daoEnt) GetRecByQN(ctx context.Context, conn tx.Conn, typeQN uniqsym.ADT) (typedef.DefRec, error) {
	qnAttr := slog.Any("typeQN", typeQN)
	rows, err := conn.Pgx.Query(ctx, selectRecByQN, uniqsym.ConvertToString(typeQN))
	if err != nil {
		dao.log.Error("query execution failed", qnAttr)
		return typedef.DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
	if err != nil {
		dao.log.Error("row scanning failed", qnAttr)
		return typedef.DefRec{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entity selection succeed", qnAttr)
	return DataToDefRec(dto)
}

func (dao *daoEnt) GetRecsByRefs(ctx context.Context, conn tx.Conn, refs []typesem.SemRef) (_ []typedef.DefRec, err error) {
	if len(refs) == 0 {
		return []typedef.DefRec{}, nil
	}
	batch := pgx.Batch{}
	for _, ref := range refs {
		if ref.TypeID.IsEmpty() {
			return nil, identity.ErrEmpty
		}
		batch.Queue(selectRecByID, ref.TypeID.String())
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []DefRec
	for _, defRef := range refs {
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", slog.Any("defRef", defRef), slog.String("q", selectRecByID))
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
		if err != nil {
			dao.log.Error("row scanning failed", slog.Any("defRef", defRef))
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return DataToDefRecs(dtos)
}

func (dao *daoEnt) SelectEnv(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (map[uniqsym.ADT]typedef.DefRec, error) {
	recs, err := dao.GetRecsByQNs(ctx, conn, typeQNs)
	if err != nil {
		return nil, err
	}
	env := make(map[uniqsym.ADT]typedef.DefRec, len(recs))
	for i, root := range recs {
		env[typeQNs[i]] = root
	}
	return env, nil
}

func (dao *daoEnt) GetRecsByQNs(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (_ []typedef.DefRec, err error) {
	if len(typeQNs) == 0 {
		return []typedef.DefRec{}, nil
	}
	batch := pgx.Batch{}
	for _, typeQN := range typeQNs {
		batch.Queue(selectRecByQN, uniqsym.ConvertToString(typeQN))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []DefRec
	for _, typeQN := range typeQNs {
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", slog.Any("typeQN", typeQN), slog.String("q", selectRecByQN))
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
		if err != nil {
			dao.log.Error("row scanning failed", slog.Any("typeQN", typeQN))
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return DataToDefRecs(dtos)
}

const (
	selectRecByQN = `
SELECT
    td.desc_id,
    td.exp_ck,
    de.desc_rn
FROM proc_type_defs td
LEFT JOIN desc_sems de
    ON de.desc_id = td.desc_id
LEFT JOIN desc_binds db
    ON db.desc_id = td.desc_id
WHERE db.desc_qn = $1
`

	selectRecByID = `
SELECT
    td.desc_id,
    td.exp_ck,
    de.desc_rn
FROM proc_type_defs td
LEFT JOIN desc_sems de
    ON de.desc_id = td.desc_id
WHERE td.desc_id = $1
`

	selectRefs = `
SELECT
    desc_id,
    def_rn
FROM proc_type_defs
`
)
