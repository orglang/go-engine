package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"

	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type daoPgx struct {
	qb  typedef.QueryBuilder
	log *slog.Logger
}

func NewDaoPgx(qb typedef.QueryBuilder, log *slog.Logger) typedef.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() typedef.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec typedef.DefRec) error {
	idAttr := slog.Any("typeID", rec.TypeRef.TypeID)
	dao.log.Log(ctx, lf.LevelTrace, "addition started", idAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", idAttr)
		return err
	}
	sql, args := dao.qb.InsertRec(dto)
	_, err = conn.Pgx.Exec(ctx, sql, args...)
	if err != nil {
		dao.log.Error("query execution failed", idAttr, slog.String("sql", sql))
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", idAttr)
	return nil
}

func (dao *daoPgx) ModifyRec(ctx context.Context, conn tx.Conn, rec typedef.DefRec) error {
	idAttr := slog.Any("typeID", rec.TypeRef.TypeID)
	dao.log.Log(ctx, lf.LevelTrace, "modification started", idAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", idAttr)
		return err
	}
	args := pgx.NamedArgs{
		"desc_id": dto.TypeID,
		"exp_ck":  dto.ExpCK,
	}
	ct, err := conn.Pgx.Exec(ctx, updateRec, args)
	if err != nil {
		dao.log.Error("query execution failed", idAttr)
		return err
	}
	if ct.RowsAffected() == 0 {
		dao.log.Error("entity update failed", idAttr)
		return typedef.ErrOptimisticUpdate(rec.TypeRef.TypeRN - 1)
	}
	dao.log.Log(ctx, lf.LevelTrace, "modification succeed", idAttr)
	return nil
}

func (dao *daoPgx) GetRefs(ctx context.Context, conn tx.Conn) ([]typesem.SemRef, error) {
	rows, err := conn.Pgx.Query(ctx, selectRefs)
	if err != nil {
		dao.log.Error("query execution failed")
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[typesem.SemRefDS])
	if err != nil {
		dao.log.Error("rows scanning failed")
		return nil, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return typesem.DataToRefs(dtos)
}

func (dao *daoPgx) GetRecByRef(ctx context.Context, conn tx.Conn, ref typesem.SemRef) (typedef.DefRec, error) {
	refAttr := slog.Any("ref", ref)
	rows, err := conn.Pgx.Query(ctx, selectRecByID, ref.TypeID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return typedef.DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[typedef.DefRecDS])
	if err != nil {
		dao.log.Error("row scanning failed", refAttr)
		return typedef.DefRec{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entity selection succeed", refAttr)
	return DataToDefRec(dto)
}

func (dao *daoPgx) GetRecByQN(ctx context.Context, conn tx.Conn, typeQN uniqsym.ADT) (typedef.DefRec, error) {
	qnAttr := slog.Any("typeQN", typeQN)
	rows, err := conn.Pgx.Query(ctx, selectRecByQN, uniqsym.ConvertToString(typeQN))
	if err != nil {
		dao.log.Error("query execution failed", qnAttr)
		return typedef.DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[typedef.DefRecDS])
	if err != nil {
		dao.log.Error("row scanning failed", qnAttr)
		return typedef.DefRec{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entity selection succeed", qnAttr)
	return DataToDefRec(dto)
}

func (dao *daoPgx) GetRecsByRefs(ctx context.Context, conn tx.Conn, refs []typesem.SemRef) (_ []typedef.DefRec, err error) {
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
	var dtos []typedef.DefRecDS
	for _, defRef := range refs {
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", slog.Any("defRef", defRef), slog.String("q", selectRecByID))
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[typedef.DefRecDS])
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

func (dao *daoPgx) SelectEnv(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (map[uniqsym.ADT]typedef.DefRec, error) {
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

func (dao *daoPgx) GetRecsByQNs(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (_ []typedef.DefRec, err error) {
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
	var dtos []typedef.DefRecDS
	for _, typeQN := range typeQNs {
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", slog.Any("typeQN", typeQN), slog.String("q", selectRecByQN))
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[typedef.DefRecDS])
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
	updateRec = `
		update proc_type_defs
		set def_rn = @def_rn,
			exp_ck = @exp_ck
		where desc_id = @desc_id
			and def_rn = @def_rn - 1`

	selectRecByQN = `
		select
			td.desc_id,
			td.exp_ck,
			de.desc_rn
		from proc_type_defs td
		left join desc_sems de
			on de.desc_id = td.desc_id
		left join desc_binds db
			on db.desc_id = td.desc_id
		where db.desc_qn = $1`

	selectRecByID = `
		select
			td.desc_id,
			td.exp_ck,
			de.desc_rn
		from proc_type_defs td
		left join desc_sems de
			on de.desc_id = td.desc_id
		where td.desc_id = $1`

	selectRefs = `
		select
			desc_id,
			def_rn
		from proc_type_defs`
)
