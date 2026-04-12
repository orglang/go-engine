package poolcomm

import (
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/uniqsym"

	"github.com/jackc/pgx/v5"
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

func (dao *pgxDAO) AddRec(source db.Source, rec ConnRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	dto := DataFromRec(rec)
	refAttr := slog.Any("ref", rec.CommRef)
	sql, args := dao.qb.insertRec(dto)
	_, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "addition succeed", slog.Any("dto", dto))
	return nil
}

func (dao *pgxDAO) GetRefsByQNs(source db.Source, qns []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error) {
	panic("unimplemented")
}

func (dao *pgxDAO) GetSnapByQry(source db.Source, qry CommQry) (CommSnap, error) {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", qry.CommRef)
	dto := DataFromQry(qry)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting started", slog.Any("qry", dto))
	sql, args := dao.qb.selectSnap(dto)
	rows, execErr := ds.Conn.Query(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return CommSnap{}, execErr
	}
	defer rows.Close()
	snapDTO, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[commSnapDS])
	if scanErr != nil {
		dao.log.Error("rows scanning failed", refAttr)
		return CommSnap{}, scanErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting succeed", slog.Any("dto", snapDTO))
	snap, convErr := DataToSnap(snapDTO)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return CommSnap{}, convErr
	}
	return snap, nil
}

func (dao *pgxDAO) ModifyRec(source db.Source, mod CommMod) error {
	if mod.CommON == nil {
		return nil
	}
	ds := db.MustConform[db.SourcePgx](source)
	dto := DataFromMod(mod)
	refAttr := slog.Any("ref", mod.CommRef)
	sql, args := dao.qb.updateRec(dto)
	_, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "modification succeed", slog.Any("dto", dto))
	return nil
}
