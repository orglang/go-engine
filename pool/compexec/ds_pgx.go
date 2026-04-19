package compexec

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/compsem"
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
	refAttr := slog.Any("ref", rec.CompRef)
	sql, args := dao.qb.insertRec(dto)
	_, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "insertion succeed", slog.Any("dto", dto))
	return nil
}

func (dao *pgxDAO) ModifyRec(db.Source, ExecMod) error {
	panic("unimplemented")
}

func (dao *pgxDAO) GetRecByRef(source db.Source, ref compsem.SemRef) (ExecRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", ref)
	sql, args := dao.qb.selectRecByRef(compsem.DataFromRef(ref))
	rows, execErr := ds.Conn.Query(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql), slog.Any("args", args))
		return ExecRec{}, execErr
	}
	defer rows.Close()
	dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[execRecDS])
	if scanErr != nil {
		dao.log.Error("rows scanning failed", refAttr, slog.Any("dto", dto))
		return ExecRec{}, scanErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dto", dto))
	rec, convErr := DataToExecRec(dto)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return ExecRec{}, convErr
	}
	return rec, nil
}

func (dao *pgxDAO) GetSnapMapByQNs(source db.Source, termQNs []uniqsym.ADT) (_ map[uniqsym.ADT]ExecSnap1, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection started", slog.Any("qns", termQNs))
	if len(termQNs) == 0 {
		return map[uniqsym.ADT]ExecSnap1{}, nil
	}
	batch := pgx.Batch{}
	for _, termQN := range termQNs {
		sql, args := dao.qb.selectSnapByQN(uniqsym.ConvertToString(termQN))
		batch.Queue(sql, args...)
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]execSnapDS, len(termQNs))
	for _, termQN := range termQNs {
		qnAttr := slog.Any("qn", termQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr)
			return nil, readErr
		}
		var dto execSnapDS
		scanErr := pgxscan.ScanOne(dto, rows)
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return nil, scanErr
		}
		dtos[termQN] = dto
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	return DataToSnapMap(dtos)
}
