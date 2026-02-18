package xactdef

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/descexec"
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
	idAttr := slog.Any("xactID", rec.XactID)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity insertion started", idAttr)
	dto, convertErr := DataFromDefRec(rec)
	if convertErr != nil {
		dao.log.Error("model conversion failed", idAttr)
		return convertErr
	}
	args := pgx.NamedArgs{
		"desc_id": dto.XactID,
		"exp_vk":  dto.ExpVK,
	}
	_, execErr := ds.Conn.Exec(ds.Ctx, insertRec, args)
	if execErr != nil {
		dao.log.Error("query execution failed", idAttr)
		return execErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity insertion succeed", idAttr)
	return nil
}

func (dao *pgxDAO) Update(source db.Source, rec DefRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	idAttr := slog.Any("xactID", rec.XactID)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity update started", idAttr)
	dto, convertErr := DataFromDefRec(rec)
	if convertErr != nil {
		dao.log.Error("model conversion failed", idAttr)
		return convertErr
	}
	args := pgx.NamedArgs{
		"desc_id": dto.XactID,
		"exp_vk":  dto.ExpVK,
	}
	_, execErr := ds.Conn.Exec(ds.Ctx, updateRec, args)
	if execErr != nil {
		dao.log.Error("query execution failed", idAttr)
		return execErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity update succeed", idAttr)
	return nil
}

func (dao *pgxDAO) SelectRefs(source db.Source) ([]descexec.ExecRef, error) {
	ds := db.MustConform[db.SourcePgx](source)
	rows, err := ds.Conn.Query(ds.Ctx, selectRefs)
	if err != nil {
		dao.log.Error("query execution failed", slog.String("q", selectRefs))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[descexec.ExecRefDS])
	if err != nil {
		dao.log.Error("rows collection failed")
		return nil, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return descexec.DataToRefs(dtos)
}

func (dao *pgxDAO) SelectRecByRef(source db.Source, ref descexec.ExecRef) (DefRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("defRef", ref)
	rows, err := ds.Conn.Query(ds.Ctx, selectRecByID, ref.DescID.String())
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
		dao.log.Error("row collection failed", qnAttr)
		return DefRec{}, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity selection succeed", qnAttr)
	return DataToDefRec(dto)
}

func (dao *pgxDAO) SelectRecsByRefs(source db.Source, refs []descexec.ExecRef) (_ []DefRec, err error) {
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

func (dao *pgxDAO) SelectRecsByQNs(source db.Source, xactQNs []uniqsym.ADT) (_ []DefRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	if len(xactQNs) == 0 {
		return []DefRec{}, nil
	}
	batch := pgx.Batch{}
	for _, xactQN := range xactQNs {
		batch.Queue(selectRecByQN, uniqsym.ConvertToString(xactQN))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []defRecDS
	for _, xactQN := range xactQNs {
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", slog.Any("xactQN", xactQN))
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
		if err != nil {
			dao.log.Error("row collection failed", slog.Any("xactQN", xactQN))
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
		insert into xact_defs (
			desc_id, exp_vk
		) values (
			@desc_id, @exp_vk
		)`

	updateRec = `
		update xact_defs
		set def_rn = @def_rn,
			exp_vk = @exp_vk
		where desc_id = @desc_id
			and def_rn = @def_rn - 1`

	selectRefs = `
		select
			desc_id,
			def_rn
		from xact_defs`

	selectRecByQN = `
		select
			xd.desc_id,
			xd.def_rn,
			xd.exp_vk
		from xact_defs xd
		left join desc_binds db
			on db.desc_id = xd.desc_id
		where db.desc_qn = $1`

	selectRecByID = `
		select
			desc_id,
			def_rn,
			exp_vk
		from xact_defs
		where desc_id = $1`
)
