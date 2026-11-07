package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/tx"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/typesem"
	"orglang/go-engine/adt/uniqsym"

	typedef "orglang/go-engine/pool/typedef/core"
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
	refAttr := slog.Any("ref", rec.TypeRef)
	dao.log.Log(ctx, lf.LevelTrace, "addition started", refAttr)
	dto, convErr := DataFromDefRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return convErr
	}
	sql, args := dao.qb.InsertRec(dto)
	_, execErr := conn.Pgx.Exec(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", refAttr)
	return nil
}

func (dao *daoPgx) Update(ctx context.Context, conn tx.Conn, rec typedef.DefRec) error {
	refAttr := slog.Any("ref", rec.TypeRef)
	dao.log.Log(ctx, lf.LevelTrace, "update started", refAttr)
	dto, convErr := DataFromDefRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return convErr
	}
	args := pgx.NamedArgs{
		"desc_id": dto.TypeID,
		"exp_vk":  dto.ExpVK,
	}
	_, execErr := conn.Pgx.Exec(ctx, updateRec, args)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr)
		return execErr
	}
	dao.log.Log(ctx, lf.LevelTrace, "update succeed", refAttr)
	return nil
}

func (dao *daoPgx) SelectRefs(ctx context.Context, conn tx.Conn) ([]typesem.SemRef, error) {
	rows, err := conn.Pgx.Query(ctx, selectRefs)
	if err != nil {
		dao.log.Error("query execution failed", slog.String("q", selectRefs))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[typesem.SemRefDS])
	if err != nil {
		dao.log.Error("rows scanning failed")
		return nil, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return typesem.DataToRefs(dtos)
}

func (dao *daoPgx) SelectRecByRef(ctx context.Context, conn tx.Conn, ref typesem.SemRef) (typedef.DefRec, error) {
	refAttr := slog.Any("defRef", ref)
	rows, err := conn.Pgx.Query(ctx, selectRecByID, ref.TypeID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("q", selectRecByID))
		return typedef.DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[typedef.DefRecDS])
	if err != nil {
		dao.log.Error("row scanning failed", refAttr)
		return typedef.DefRec{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", refAttr)
	return DataToDefRec(dto)
}

func (dao *daoPgx) SelectRecByQN(ctx context.Context, conn tx.Conn, xactQN uniqsym.ADT) (typedef.DefRec, error) {
	qnAttr := slog.Any("xactQN", xactQN)
	rows, err := conn.Pgx.Query(ctx, selectRecByQN, uniqsym.ConvertToString(xactQN))
	if err != nil {
		dao.log.Error("query execution failed", qnAttr, slog.String("q", selectRecByQN))
		return typedef.DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[typedef.DefRecDS])
	if err != nil {
		dao.log.Error("row scanning failed", qnAttr)
		return typedef.DefRec{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", qnAttr)
	return DataToDefRec(dto)
}

func (dao *daoPgx) SelectRecsByRefs(ctx context.Context, conn tx.Conn, refs []typesem.SemRef) (_ []typedef.DefRec, err error) {
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
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return DataToDefRecs(dtos)
}

func (dao *daoPgx) GetRecsByQNs(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (_ map[uniqsym.ADT]typedef.DefRec, err error) {
	if len(typeQNs) == 0 {
		return map[uniqsym.ADT]typedef.DefRec{}, nil
	}
	batch := pgx.Batch{}
	sql := dao.qb.SelectRecByQN()
	for _, typeQN := range typeQNs {
		batch.Queue(sql, uniqsym.ConvertToString(typeQN))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]typedef.DefRecDS, len(typeQNs))
	for _, typeQN := range typeQNs {
		qnAttr := slog.Any("qn", typeQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr, slog.Any("sql", sql))
			return map[uniqsym.ADT]typedef.DefRec{}, readErr
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[typedef.DefRecDS])
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return map[uniqsym.ADT]typedef.DefRec{}, scanErr
		}
		dtos[typeQN] = dto
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return DataToDefRecMap(dtos)
}

const (
	updateRec = `
		update pool_type_defs
		set def_rn = @def_rn,
			exp_vk = @exp_vk
		where desc_id = @desc_id
			and def_rn = @def_rn - 1`

	selectRefs = `
		select
			desc_id,
			def_rn
		from pool_type_defs`

	selectRecByQN = `
		select
			xd.desc_id,
			xd.def_rn,
			xd.exp_vk
		from pool_type_defs xd
		left join desc_binds db
			on db.desc_id = xd.desc_id
		where db.desc_qn = $1`

	selectRecByID = `
		select
			desc_id,
			def_rn,
			exp_vk
		from pool_type_defs
		where desc_id = $1`
)
