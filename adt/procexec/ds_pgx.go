package procexec

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/procstep"
	"orglang/go-engine/adt/revnum"
)

// Adapter
type pgxDAO struct {
	log *slog.Logger
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func newPgxDAO(l *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{l.With(name)}
}

func (dao *pgxDAO) InsertRec(source db.Source, rec ExecRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	dto := DataFromExecRec(rec)
	args := pgx.NamedArgs{
		"impl_id": dto.ImplID,
		"impl_rn": dto.ImplRN,
		"chnl_ph": dto.ChnlPH,
	}
	refAttr := slog.Any("ref", rec.ImplRef)
	_, execErr := ds.Conn.Exec(ds.Ctx, insertRec, args)
	if execErr != nil {
		dao.log.Error("execution failed", refAttr)
		return execErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "insertion succeed", slog.Any("dto", dto))
	return nil
}

func (dao *pgxDAO) SelectSnap(source db.Source, ref implsem.SemRef) (ExecSnap, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", ref)
	chnlRows, err := ds.Conn.Query(ds.Ctx, selectChnls, ref.ImplID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return ExecSnap{}, err
	}
	defer chnlRows.Close()
	chnlDtos, err := pgx.CollectRows(chnlRows, pgx.RowToStructByName[implvar.VarRecDS])
	if err != nil {
		dao.log.Error("row collection failed", refAttr, slog.Any("t", reflect.TypeOf(chnlDtos)))
		return ExecSnap{}, err
	}
	chnls, err := implvar.DataToVarRecs(chnlDtos)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return ExecSnap{}, err
	}
	stepRows, err := ds.Conn.Query(ds.Ctx, selectSteps, ref.ImplID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return ExecSnap{}, err
	}
	defer stepRows.Close()
	stepDtos, err := pgx.CollectRows(stepRows, pgx.RowToStructByName[procstep.StepRecDS])
	if err != nil {
		dao.log.Error("row collection failed", refAttr, slog.Any("t", reflect.TypeOf(stepDtos)))
		return ExecSnap{}, err
	}
	steps, err := procstep.DataToStepRecs(stepDtos)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return ExecSnap{}, err
	}
	dao.log.Debug("snap selection succeed", refAttr)
	return ExecSnap{
		ChnlVRs: implvar.IndexBy(ChnlPH, chnls),
		ProcSRs: implvar.IndexBy(procstep.ChnlID, steps),
	}, nil
}

func (dao *pgxDAO) UpdateProc(source db.Source, mod ExecMod) (err error) {
	if len(mod.Locks) == 0 {
		panic("empty locks")
	}
	ds := db.MustConform[db.SourcePgx](source)
	dto, err := DataFromMod(mod)
	if err != nil {
		dao.log.Error("conversion failed")
		return err
	}
	// binds
	bindReq := pgx.Batch{}
	for _, dto := range dto.Binds {
		args := pgx.NamedArgs{
			"impl_id":  dto.ImplID,
			"impl_rn":  dto.ImplRN,
			"chnl_ph":  dto.ChnlPH,
			"chnl_id":  dto.ChnlID,
			"state_id": dto.ExpVK,
		}
		bindReq.Queue(insertBind, args)
	}
	if bindReq.Len() > 0 {
		bindRes := ds.Conn.SendBatch(ds.Ctx, &bindReq)
		defer func() {
			err = errors.Join(err, bindRes.Close())
		}()
		for _, dto := range dto.Binds {
			_, err = bindRes.Exec()
			if err != nil {
				dao.log.Error("execution failed", slog.Any("dto", dto))
			}
		}
		if err != nil {
			return err
		}
	}
	// steps
	stepReq := pgx.Batch{}
	for _, dto := range dto.Steps {
		args := pgx.NamedArgs{
			"kind":    dto.K,
			"proc_id": dto.ExecID,
			"chnl_id": dto.ChnlID,
			"proc_er": dto.ProcER,
		}
		stepReq.Queue(insertStep, args)
	}
	if stepReq.Len() > 0 {
		stepRes := ds.Conn.SendBatch(ds.Ctx, &stepReq)
		defer func() {
			err = errors.Join(err, stepRes.Close())
		}()
		for _, dto := range dto.Steps {
			_, err = stepRes.Exec()
			if err != nil {
				dao.log.Error("execution failed", slog.Any("dto", dto))
			}
		}
		if err != nil {
			return err
		}
	}
	// execs
	execReq := pgx.Batch{}
	for _, dto := range dto.Locks {
		args := pgx.NamedArgs{
			"impl_id": dto.ImplID,
			"impl_rn": dto.ImplRN,
		}
		execReq.Queue(updateExec, args)
	}
	execRes := ds.Conn.SendBatch(ds.Ctx, &execReq)
	defer func() {
		err = errors.Join(err, execRes.Close())
	}()
	for _, dto := range dto.Locks {
		ct, err := execRes.Exec()
		if err != nil {
			dao.log.Error("execution failed", slog.Any("dto", dto))
		}
		if ct.RowsAffected() == 0 {
			dao.log.Error("update failed")
			return errOptimisticUpdate(revnum.ADT(dto.ImplRN))
		}
	}
	if err != nil {
		return err
	}
	dao.log.Debug("update succeed")
	return nil
}

const (
	insertRec = `
		insert into proc_execs (
			impl_id, impl_rn, chnl_ph
		) values (
			@impl_id, @impl_rn, @chnl_ph
		)`

	insertBind = `
		insert into proc_binds (
			impl_id, chnl_ph, chnl_id, state_id, impl_rn
		) values (
			@impl_id, @chnl_ph, @chnl_id, @state_id, @impl_rn
		)`

	insertStep = `
		insert into proc_steps (
			impl_id, chnl_id, kind, proc_er
		) values (
			@impl_id, @chnl_id, @kind, @proc_er
		)`

	updateExec = `
		update proc_execs
		set impl_rn = @impl_rn + 1
		where impl_id = @impl_id
			and impl_rn = @impl_rn`

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
		left join pool_execs prvd
			on prvd.desc_id = liab.desc_id`

	selectSteps = ``
)
