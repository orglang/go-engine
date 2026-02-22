package poolexec

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/uniqsym"
)

type pgxDAO struct {
	log *slog.Logger
}

func newPgxDAO(log *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{log.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) InsertRec(source db.Source, rec ExecRec) (err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dto := DataFromExecRec(rec)
	args := pgx.NamedArgs{
		"chnl_id": dto.ChnlID,
		"chnl_ph": dto.ChnlPH,
		"exp_vk":  dto.ExpVK,
		"impl_id": dto.ImplID,
		"impl_rn": dto.ImplRN,
	}
	refAttr := slog.Any("ref", rec.ImplRef)
	_, err = ds.Conn.Exec(ds.Ctx, insertRec, args)
	if err != nil {
		dao.log.Error("execution failed", refAttr)
		return err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "insertion succeed", refAttr)
	return nil
}

func (dao *pgxDAO) TouchRec(source db.Source, ref implsem.SemRef) error {
	return nil
}

func (dao *pgxDAO) TouchRecs(db.Source, []implsem.SemRef) error {
	return nil
}

func (dao *pgxDAO) SelectRecsByQNs(source db.Source, implQNs []uniqsym.ADT) (_ map[uniqsym.ADT]ExecRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "starting selection...", slog.Any("qns", implQNs))
	if len(implQNs) == 0 {
		return map[uniqsym.ADT]ExecRec{}, nil
	}
	batch := pgx.Batch{}
	for _, implQN := range implQNs {
		batch.Queue(selectRecByQN, uniqsym.ConvertToString(implQN))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]execRecDS, len(implQNs))
	for _, implQN := range implQNs {
		qnAttr := slog.Any("qn", implQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr)
			return nil, readErr
		}
		dto, collectErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[execRecDS])
		if collectErr != nil {
			dao.log.Error("row collection failed", qnAttr)
			return nil, collectErr
		}
		dtos[implQN] = dto
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	return DataToRefMap(dtos)
}

func (dao *pgxDAO) SelectSubs(source db.Source, ref implsem.SemRef) (ExecSnap, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("execRef", ref)
	rows, err := ds.Conn.Query(ds.Ctx, selectOrgSnap, ref.ImplID.String())
	if err != nil {
		dao.log.Error("execution failed", refAttr)
		return ExecSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[execSnapDS])
	if err != nil {
		dao.log.Error("collection failed", refAttr, slog.Any("type", reflect.TypeFor[execSnapDS]))
		return ExecSnap{}, err
	}
	snap, err := DataToExecSnap(dto)
	if err != nil {
		dao.log.Error("conversion failed")
		return ExecSnap{}, err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", refAttr)
	return snap, nil
}

func (dao *pgxDAO) SelectRefs(source db.Source) ([]implsem.SemRef, error) {
	ds := db.MustConform[db.SourcePgx](source)
	query := `
		select
			impl_id, impl_rn
		from pool_execs`
	rows, err := ds.Conn.Query(ds.Ctx, query)
	if err != nil {
		dao.log.Error("execution failed", slog.String("q", query))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[implsem.SemRefDS])
	if err != nil {
		dao.log.Error("collection failed", slog.Any("type", reflect.TypeFor[[]implsem.SemRefDS]))
		return nil, err
	}
	refs, err := implsem.DataToRefs(dtos)
	if err != nil {
		dao.log.Error("conversion failed")
		return nil, err
	}
	return refs, nil
}

const (
	insertRec = `
		insert into pool_execs (
			impl_id, impl_rn, chnl_id, chnl_ph, exp_vk
		) values (
			@impl_id, @impl_rn, @chnl_id, @chnl_ph, @exp_vk
		)`

	selectRecByQN = `
		select
			is.impl_id,
			is.impl_rn,
			pe.chnl_id,
			pe.chnl_ph,
			pe.exp_vk
		from pool_execs pe
		left join impl_sems is
			on is.impl_id = pe.impl_id
		left join impl_binds ib
			on ib.impl_id = pe.impl_id
		where ib.impl_qn = $1`

	insertStep = `
		insert into pool_steps (
			proc_id, chnl_id, kind, spec
		) values (
			@proc_id, @chnl_id, @kind, @spec
		)`

	updateRoot = `
		update pool_roots
		set rev = @rev + 1
		where desc_id = @desc_id
			and rev = @rev`

	selectOrgSnap = `
		select
			sup.desc_id,
			sup.title,
			jsonb_agg(json_build_object('desc_id', sub.desc_id, 'title', sub.title)) as subs
		from pool_roots sup
		left join pool_sups rel
			on rel.sup_pool_id = sup.desc_id
		left join pool_roots sub
			on sub.desc_id = rel.desc_id
			and sub.rev = rel.rev
		where sup.desc_id = $1
		group by sup.desc_id, sup.title`

	selectChnls = `
		with bnds as not materialized (
			select distinct on (chnl_ph)
				*
			from proc_bnds
			where proc_id = 'proc1'
			order by chnl_ph, abs(rev) desc
		), liabs as not materialized (
			select distinct on (proc_id)
				*
			from pool_liabs
			where proc_id = 'proc1'
			order by proc_id, abs(rev) desc
		)
		select
			bnd.*,
			prvd.desc_id
		from bnds bnd
		left join liabs liab
			on liab.proc_id = bnd.proc_id
		left join pool_roots prvd
			on prvd.desc_id = liab.desc_id`
)
