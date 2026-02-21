package poolexec

import (
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"
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
		"exec_id": dto.ImplID,
		"exec_rn": dto.ImplRN,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertExec, args)
	if err != nil {
		dao.log.Error("execution failed")
		return err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "insertion succeed", slog.Any("execRef", rec.ExecRef))
	return nil
}

func (dao *pgxDAO) InsertLiab(source db.Source, liab Liab) (err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dto := DataFromLiab(liab)
	args := pgx.NamedArgs{
		"exec_id": dto.ImplID,
		"exec_rn": dto.ImplRN,
		"proc_id": dto.ProcID,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertLiab, args)
	if err != nil {
		dao.log.Error("execution failed")
		return err
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "insertion succeed", slog.Any("execRef", liab.ExecRef))
	return nil
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
			exec_id, exec_rn
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
	insertExec = `
		insert into pool_execs (
			exec_id, exec_rn, proc_id
		) values (
			@exec_id, @exec_rn, @proc_id
		)`

	insertLiab = `
		insert into pool_liabs (
			desc_id, proc_id, rev
		) values (
			@desc_id, @proc_id, @rev
		)`

	insertBind = `
		insert into pool_assets (
			desc_id, chnl_key, state_id, ex_pool_id, rev
		) values (
			@desc_id, @chnl_key, @state_id, @ex_pool_id, @rev
		)`

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

	selectSteps = ``
)
