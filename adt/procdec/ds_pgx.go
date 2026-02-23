package procdec

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/identity"
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

func (dao *pgxDAO) InsertRec(source db.Source, rec DecRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", rec.DescRef)
	dto, err := DataFromDecRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}
	args := pgx.NamedArgs{
		"desc_id":     dto.DescID,
		"provider_vr": dto.ProviderVR,
		"client_vrs":  dto.ClientVRs,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRec, args)
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}
	return nil
}

func (dao *pgxDAO) SelectSnap(source db.Source, ref descsem.SemRef) (DecSnap, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", ref)
	rows, err := ds.Conn.Query(ds.Ctx, selectByRef, ref.DescID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return DecSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[decSnapDS])
	if err != nil {
		dao.log.Error("row collection failed", refAttr)
		return DecSnap{}, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entitiy selection succeed", slog.Any("dto", dto))
	return DataToDecSnap(dto)
}

func (dao *pgxDAO) SelectEnv(source db.Source, ids []identity.ADT) (map[identity.ADT]DecRec, error) {
	decs, err := dao.SelectRecs(source, ids)
	if err != nil {
		return nil, err
	}
	env := make(map[identity.ADT]DecRec, len(decs))
	for _, dec := range decs {
		env[dec.DescRef.DescID] = dec
	}
	return env, nil
}

func (dao *pgxDAO) SelectRecs(source db.Source, ids []identity.ADT) (_ []DecRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	if len(ids) == 0 {
		return []DecRec{}, nil
	}
	batch := pgx.Batch{}
	for _, rid := range ids {
		if rid.IsEmpty() {
			return nil, identity.ErrEmpty
		}
		batch.Queue(selectByRef, rid.String())
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []decRecDS
	for _, rid := range ids {
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", slog.Any("id", rid), slog.String("q", selectByRef))
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[decRecDS])
		if err != nil {
			dao.log.Error("row collection failed", slog.Any("id", rid))
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return DataToDecRecs(dtos)
}

func (dao *pgxDAO) SelectRefs(source db.Source) ([]descsem.SemRef, error) {
	ds := db.MustConform[db.SourcePgx](source)
	query := `
		select
			desc_id, rev, title
		from dec_roots`
	rows, err := ds.Conn.Query(ds.Ctx, query)
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
	return descsem.DataToRefs(dtos)
}

const (
	insertRec = `
		insert into proc_decs (
			desc_id, provider_vr, client_vrs
		) values (
			@desc_id, @provider_vr, @client_vrs
		)`

	selectByRef = `
		select
			pd.desc_id,
			ds.desc_rn,
			pd.provider_vr,
			pd.client_vrs
		from proc_decs pd
		left join desc_sems ds
			on ds.desc_id = pd.desc_id
		where pd.desc_id = $1`
)
