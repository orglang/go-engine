package termdef

import (
	"log/slog"
	"reflect"

	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/lib/db"

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

func (dao *pgxDAO) AddRec(source db.Source, rec DefRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", rec.TermRef)
	dto, convErr := DataFromDecRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return convErr
	}
	sql, args := dao.qb.insertRec(dto)
	_, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr)
		return execErr
	}
	return nil
}

func (dao *pgxDAO) GetRecByQN(source db.Source, qn uniqsym.ADT) (DefRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	qnAttr := slog.Any("qn", qn)
	sql, args := dao.qb.selectRecByQN(uniqsym.ConvertToString(qn))
	rows, execErr := ds.Conn.Query(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", qnAttr)
		return DefRec{}, execErr
	}
	dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
	if scanErr != nil {
		dao.log.Error("rows scanning failed", qnAttr)
		return DefRec{}, scanErr
	}
	rec, convErr := DataToDecRec(dto)
	if convErr != nil {
		dao.log.Error("model conversion failed", qnAttr)
		return DefRec{}, convErr
	}
	return rec, nil
}
