package typedef

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

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
	refAttr := slog.Any("defRef", rec.DefRef)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity insertion started", refAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}
	defArgs := pgx.NamedArgs{
		"def_id": dto.ID,
		"def_rn": dto.RN,
		"exp_id": dto.ExpID,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRec, defArgs)
	if err != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("q", insertRec))
		return err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity insertion succeed", refAttr)
	return nil
}

func (dao *pgxDAO) Update(source db.Source, rec DefRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("defRef", rec.DefRef)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity update started", refAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}
	args := pgx.NamedArgs{
		"def_id": dto.ID,
		"def_rn": dto.RN,
		"exp_id": dto.ExpID,
	}
	ct, err := ds.Conn.Exec(ds.Ctx, updateRec, args)
	if err != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("q", updateRec))
		return err
	}
	if ct.RowsAffected() == 0 {
		dao.log.Error("entity update failed", refAttr)
		return errOptimisticUpdate(rec.DefRef.RN - 1)
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity update succeed", refAttr)
	return nil
}

func (dao *pgxDAO) SelectRefs(source db.Source) ([]DefRef, error) {
	ds := db.MustConform[db.SourcePgx](source)
	rows, err := ds.Conn.Query(ds.Ctx, selectRefs)
	if err != nil {
		dao.log.Error("query execution failed", slog.String("q", selectRefs))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[defRefDS])
	if err != nil {
		dao.log.Error("rows collection failed")
		return nil, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return DataToDefRefs(dtos)
}

func (dao *pgxDAO) SelectRecByRef(source db.Source, defRef DefRef) (DefRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("defRef", defRef)
	rows, err := ds.Conn.Query(ds.Ctx, selectRecByID, defRef.ID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("q", selectRecByID))
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
		dao.log.Error("query execution failed", qnAttr, slog.String("q", selectRecByQN))
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

func (dao *pgxDAO) SelectRecsByRefs(source db.Source, defRefs []DefRef) (_ []DefRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	if len(defRefs) == 0 {
		return []DefRec{}, nil
	}
	batch := pgx.Batch{}
	for _, defRef := range defRefs {
		if defRef.ID.IsEmpty() {
			return nil, identity.ErrEmpty
		}
		batch.Queue(selectRecByID, defRef.ID.String())
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []defRecDS
	for _, defRef := range defRefs {
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
			def_id, def_rn, exp_id
		) values (
			@def_id, @def_rn, @exp_id
		)`

	updateRec = `
		update type_defs
		set def_rn = @def_rn,
			exp_id = @exp_id
		where def_id = @def_id
			and def_rn = @def_rn - 1`

	selectRecByQN = `
		select
			td.def_id,
			td.def_rn,
			td.exp_id
		from type_defs td
		left join syn_decs sd
			on sd.dec_id = td.def_id
			and sd.from_rn >= td.def_rn
			and sd.to_rn > td.def_rn
		where sd.dec_qn = $1`

	selectRecByID = `
		select
			def_id,
			def_rn,
			exp_id
		from type_defs
		where def_id = $1`

	selectRefs = `
		select
			def_id,
			def_rn
		from type_defs`
)
