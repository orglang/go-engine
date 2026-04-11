package poolstep

import (
	"errors"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

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

func (dao *pgxDAO) AddRec(db.Source, StepRec) error {
	panic("unimplemented")
}

func (dao *pgxDAO) AddRecs(source db.Source, recs []StepRec) (err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelDebug, "insertion started", slog.Any("recs", recs))
	batch := pgx.Batch{}
	for _, rec := range recs {
		dto := DataFromStepRec(rec)
		sql, args := dao.qb.insertRec(dto)
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
