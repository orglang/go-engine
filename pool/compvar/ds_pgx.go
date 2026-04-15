package compvar

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/compvar"
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

func (dao *pgxDAO) AddRecs(source db.Source, recs []compvar.VarRec) (err error) {
	ds := db.MustConform[db.SourcePgx](source)
	batch := pgx.Batch{}
	for _, rec := range recs {
		dto := compvar.DataFromVarRec(rec)
		sql, args := dao.qb.insertRec(getTableName(rec), dto)
		batch.Queue(sql, args...)
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for _, rec := range recs {
		_, readErr := br.Exec()
		if readErr != nil {
			dao.log.Error("query execution failed", slog.Any("rec", rec))
			return readErr
		}
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "insertion succeed")
	return nil
}

func getTableName(rec compvar.VarRec) string {
	switch rec.(type) {
	case compvar.StructRec:
		return poolStructVars
	case compvar.LinearRec:
		return poolLinearVars
	default:
		panic(compvar.ErrUnexpectedRecType(rec))
	}
}
