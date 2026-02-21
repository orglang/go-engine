package typedef

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/uniqsym"
)

type pgxDAO struct {
	log *slog.Logger
}

func newPgxDAO(l *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{l.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) InsertRec(source db.Source, rec DefRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	idAttr := slog.Any("typeID", rec.DescRef.DescID)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity insertion started", idAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", idAttr)
		return err
	}
	args := pgx.NamedArgs{
		"desc_id": dto.DescID,
		"exp_vk":  dto.ExpVK,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRec, args)
	if err != nil {
		dao.log.Error("query execution failed", idAttr)
		return err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity insertion succeed", idAttr)
	return nil
}

func (dao *pgxDAO) Update(source db.Source, rec DefRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	idAttr := slog.Any("typeID", rec.DescRef.DescID)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity update started", idAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", idAttr)
		return err
	}
	args := pgx.NamedArgs{
		"desc_id": dto.DescID,
		"exp_vk":  dto.ExpVK,
	}
	ct, err := ds.Conn.Exec(ds.Ctx, updateRec, args)
	if err != nil {
		dao.log.Error("query execution failed", idAttr)
		return err
	}
	if ct.RowsAffected() == 0 {
		dao.log.Error("entity update failed", idAttr)
		return errOptimisticUpdate(rec.DescRef.DescRN - 1)
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity update succeed", idAttr)
	return nil
}

func (dao *pgxDAO) SelectRefs(source db.Source) ([]descsem.SemRef, error) {
	ds := db.MustConform[db.SourcePgx](source)
	rows, err := ds.Conn.Query(ds.Ctx, selectRefs)
	if err != nil {
		dao.log.Error("query execution failed")
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[descsem.SemRefDS])
	if err != nil {
		dao.log.Error("rows collection failed")
		return nil, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return descsem.DataToRefs(dtos)
}

func (dao *pgxDAO) SelectRecByRef(source db.Source, ref descsem.SemRef) (DefRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", ref)
	rows, err := ds.Conn.Query(ds.Ctx, selectRecByID, ref.DescID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
	if err != nil {
		dao.log.Error("row collection failed", refAttr)
		return DefRec{}, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity selection succeed", refAttr)
	return DataToDefRec(dto)
}

func (dao *pgxDAO) SelectRecByQN(source db.Source, typeQN uniqsym.ADT) (DefRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	qnAttr := slog.Any("typeQN", typeQN)
	rows, err := ds.Conn.Query(ds.Ctx, selectRecByQN, uniqsym.ConvertToString(typeQN))
	if err != nil {
		dao.log.Error("query execution failed", qnAttr)
		return DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
	if err != nil {
		dao.log.Error("row collection failed", qnAttr)
		return DefRec{}, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity selection succeed", qnAttr)
	return DataToDefRec(dto)
}

func (dao *pgxDAO) SelectRecsByRefs(source db.Source, refs []descsem.SemRef) (_ []DefRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	if len(refs) == 0 {
		return []DefRec{}, nil
	}
	batch := pgx.Batch{}
	for _, ref := range refs {
		if ref.DescID.IsEmpty() {
			return nil, identity.ErrEmpty
		}
		batch.Queue(selectRecByID, ref.DescID.String())
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
			dao.log.Error("row collection failed", slog.Any("defRef", defRef))
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return DataToDefRecs(dtos)
}

func (dao *pgxDAO) SelectEnv(source db.Source, typeQNs []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error) {
	recs, err := dao.SelectRecsByQNs(source, typeQNs)
	if err != nil {
		return nil, err
	}
	env := make(map[uniqsym.ADT]DefRec, len(recs))
	for i, root := range recs {
		env[typeQNs[i]] = root
	}
	return env, nil
}

func (dao *pgxDAO) SelectRecsByQNs(source db.Source, typeQNs []uniqsym.ADT) (_ []DefRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	if len(typeQNs) == 0 {
		return []DefRec{}, nil
	}
	batch := pgx.Batch{}
	for _, typeQN := range typeQNs {
		batch.Queue(selectRecByQN, uniqsym.ConvertToString(typeQN))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []defRecDS
	for _, typeQN := range typeQNs {
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", slog.Any("typeQN", typeQN), slog.String("q", selectRecByQN))
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
		if err != nil {
			dao.log.Error("row collection failed", slog.Any("typeQN", typeQN))
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return DataToDefRecs(dtos)
}

const (
	insertRec = `
		insert into type_defs (
			desc_id, exp_vk
		) values (
			@desc_id, @exp_vk
		)`

	updateRec = `
		update type_defs
		set def_rn = @def_rn,
			exp_vk = @exp_vk
		where desc_id = @desc_id
			and def_rn = @def_rn - 1`

	selectRecByQN = `
		select
			td.desc_id,
			td.exp_vk,
			de.desc_rn
		from type_defs td
		left join desc_sems de
			on de.desc_id = td.desc_id
		left join desc_binds db
			on db.desc_id = td.desc_id
		where db.desc_qn = $1`

	selectRecByID = `
		select
			td.desc_id,
			td.exp_vk,
			de.desc_rn
		from type_defs td
		left join desc_sems de
			on de.desc_id = td.desc_id
		where td.desc_id = $1`

	selectRefs = `
		select
			desc_id,
			def_rn
		from type_defs`
)
