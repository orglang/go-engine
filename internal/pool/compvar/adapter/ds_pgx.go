package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/lf"
	"orglang/go-engine/lib/tx"

	adtcompvar "orglang/go-engine/adt/compvar"

	compvar "orglang/go-engine/pool/compvar/core"
)

type daoPgx struct {
	qb  compvar.QueryBuilder
	log *slog.Logger
}

func NewDaoPgx(qb compvar.QueryBuilder, log *slog.Logger) compvar.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() compvar.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRecs(ctx context.Context, conn tx.Conn, recs []adtcompvar.VarRec) (err error) {
	batch := pgx.Batch{}
	for _, rec := range recs {
		dto := adtcompvar.DataFromVarRec(rec)
		sql, args := dao.qb.InsertRec(getTableName(rec), dto)
		batch.Queue(sql, args...)
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
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
	dao.log.Log(ctx, lf.LevelTrace, "insertion succeed")
	return nil
}

func getTableName(rec adtcompvar.VarRec) string {
	switch rec.(type) {
	case adtcompvar.StructRec:
		return poolStructVars
	case adtcompvar.LinearRec:
		return poolLinearVars
	default:
		panic(adtcompvar.ErrUnexpectedRecType(rec))
	}
}
