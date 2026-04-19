package typedef

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/typesem"
	"orglang/go-engine/adt/uniqsym"
)

type pgxDAO struct {
	qb  queryBuilder
	log *slog.Logger
}

func newPgxDAO(qb queryBuilder, log *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{qb, log.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) AddRec(source db.Source, rec DefRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", rec.TypeRef)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "addition started", refAttr)
	dto, convErr := DataFromDefRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return convErr
	}
	sql, args := dao.qb.insertRec(dto)
	_, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "addition succeed", refAttr)
	return nil
}

func (dao *pgxDAO) Update(source db.Source, rec DefRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", rec.TypeRef)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "update started", refAttr)
	dto, convErr := DataFromDefRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return convErr
	}
	args := pgx.NamedArgs{
		"desc_id": dto.TypeID,
		"exp_vk":  dto.ExpVK,
	}
	_, execErr := ds.Conn.Exec(ds.Ctx, updateRec, args)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr)
		return execErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "update succeed", refAttr)
	return nil
}

func (dao *pgxDAO) SelectRefs(source db.Source) ([]typesem.SemRef, error) {
	ds := db.MustConform[db.SourcePgx](source)
	rows, err := ds.Conn.Query(ds.Ctx, selectRefs)
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
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return typesem.DataToRefs(dtos)
}

func (dao *pgxDAO) SelectRecByRef(source db.Source, ref typesem.SemRef) (DefRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("defRef", ref)
	rows, err := ds.Conn.Query(ds.Ctx, selectRecByID, ref.TypeID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("q", selectRecByID))
		return DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
	if err != nil {
		dao.log.Error("row scanning failed", refAttr)
		return DefRec{}, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting succeed", refAttr)
	return DataToDefRec(dto)
}

func (dao *pgxDAO) SelectRecByQN(source db.Source, xactQN uniqsym.ADT) (DefRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	qnAttr := slog.Any("xactQN", xactQN)
	rows, err := ds.Conn.Query(ds.Ctx, selectRecByQN, uniqsym.ConvertToString(xactQN))
	if err != nil {
		dao.log.Error("query execution failed", qnAttr, slog.String("q", selectRecByQN))
		return DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
	if err != nil {
		dao.log.Error("row scanning failed", qnAttr)
		return DefRec{}, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting succeed", qnAttr)
	return DataToDefRec(dto)
}

func (dao *pgxDAO) SelectRecsByRefs(source db.Source, refs []typesem.SemRef) (_ []DefRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	if len(refs) == 0 {
		return []DefRec{}, nil
	}
	batch := pgx.Batch{}
	for _, ref := range refs {
		if ref.TypeID.IsEmpty() {
			return nil, identity.ErrEmpty
		}
		batch.Queue(selectRecByID, ref.TypeID.String())
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []defRecDS
	for _, defRef := range refs {
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", slog.Any("defRef", defRef), slog.String("q", selectRecByID))
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
		if err != nil {
			dao.log.Error("row scanning failed", slog.Any("defRef", defRef))
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return DataToDefRecs(dtos)
}

func (dao *pgxDAO) GetRecsByQNs(source db.Source, typeQNs []uniqsym.ADT) (_ map[uniqsym.ADT]DefRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	if len(typeQNs) == 0 {
		return map[uniqsym.ADT]DefRec{}, nil
	}
	batch := pgx.Batch{}
	sql := dao.qb.selectRecByQN()
	for _, typeQN := range typeQNs {
		batch.Queue(sql, uniqsym.ConvertToString(typeQN))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]defRecDS, len(typeQNs))
	for _, typeQN := range typeQNs {
		qnAttr := slog.Any("qn", typeQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr, slog.Any("sql", sql))
			return map[uniqsym.ADT]DefRec{}, readErr
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return map[uniqsym.ADT]DefRec{}, scanErr
		}
		dtos[typeQN] = dto
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
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
