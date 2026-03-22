package poolconn

import (
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/commsem"
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

func (dao *pgxDAO) AddRec(source db.Source, rec ConnRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	dto := DataFromRec(rec)
	commAttr := slog.Any("comm", rec.CommRef)
	sql, args := dao.qb.insertRec(dto)
	_, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", commAttr)
		return execErr
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "insertion succeed", slog.Any("dto", dto))
	return nil
}

func (dao *pgxDAO) GetRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error) {
	panic("unimplemented")
}

func (dao *pgxDAO) GetSnapByQry(db.Source, ConnQry) (ConnSnap, error) {
	panic("unimplemented")
}

func (dao *pgxDAO) UpdateRec(db.Source, ConnMod) error {
	panic("unimplemented")
}
