package poolexec

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"github.com/georgysavva/scany/v2/pgxscan"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/implsem"
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

func (dao *pgxDAO) AddRec(source db.Source, rec ExecRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	dto := DataFromExecRec(rec)
	implAttr := slog.Any("impl", rec.ImplRef)
	sql, args := dao.qb.insertRec(dto)
	_, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", implAttr)
		return execErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "insertion succeed", slog.Any("dto", dto))
	return nil
}

func (dao *pgxDAO) GetRecsByQNs(source db.Source, implQNs []uniqsym.ADT) (_ map[uniqsym.ADT]ExecRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection started", slog.Any("qns", implQNs))
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
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[execRecDS])
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return nil, scanErr
		}
		dtos[implQN] = dto
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	return DataToRefMap(dtos)
}

func (dao *pgxDAO) GetRefs(source db.Source) ([]implsem.SemRef, error) {
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

func (dao *pgxDAO) GetSnap(source db.Source, ref implsem.SemRef) (ExecCtxSnap, error) {
	ds := db.MustConform[db.SourcePgx](source)
	implAttr := slog.Any("impl", ref)
	refDTO, convErr1 := implsem.DataFromRef(ref)
	if convErr1 != nil {
		dao.log.Error("model converison failed", implAttr)
		return ExecCtxSnap{}, convErr1
	}
	sql, args := dao.qb.selectSnap(refDTO)
	rows, execErr := ds.Conn.Query(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", implAttr, slog.String("sql", sql), slog.Any("args", args))
		return ExecCtxSnap{}, execErr
	}
	defer rows.Close()
	snapDTO, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[execCtxSnapDS])
	if scanErr != nil {
		dao.log.Error("row scanning failed", implAttr, slog.Any("dto", snapDTO))
		return ExecCtxSnap{}, scanErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dto", snapDTO))
	snap, convErr2 := DataToExecSnap(snapDTO)
	if convErr2 != nil {
		dao.log.Error("model converison failed", implAttr)
		return ExecCtxSnap{}, convErr2
	}
	return snap, nil
}

func (dao *pgxDAO) GetSnapsByQNs(source db.Source, implQNs []uniqsym.ADT) (_ map[uniqsym.ADT]ExecLiabSnap, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection started", slog.Any("qns", implQNs))
	if len(implQNs) == 0 {
		return map[uniqsym.ADT]ExecLiabSnap{}, nil
	}
	batch := pgx.Batch{}
	for _, implQN := range implQNs {
		sql, args := dao.qb.selectSnapByQN(uniqsym.ConvertToString(implQN))
		batch.Queue(sql, args...)
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]execLiabSnapDS, len(implQNs))
	for _, implQN := range implQNs {
		qnAttr := slog.Any("qn", implQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr)
			return nil, readErr
		}
		var dto execLiabSnapDS
		scanErr := pgxscan.ScanOne(dto, rows)
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return nil, scanErr
		}
		dtos[implQN] = dto
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	return DataToSnapMap(dtos)
}

const (
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
