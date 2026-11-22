package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	adtcompvar "orglang/go-engine/internal/adt/compvar"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"

	compvar "orglang/go-engine/internal/pool/compvar/core"
)

type daoPgx struct {
	log *slog.Logger
}

func NewDaoPgx(log *slog.Logger) compvar.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{log.With(name)}
}

// for compilation purposes
func newRepo() compvar.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRecs(ctx context.Context, conn tx.Conn, recs []adtcompvar.VarRec) (err error) {
	batch := pgx.Batch{}
	for _, rec := range recs {
		dto := compvaradapter.DataFromVarRec(rec)
		batch.Queue(insertRec(getTableName(rec)), dto.CompID, dto.CompRN, dto.ConnID, dto.ChnlID, dto.ChnlPH, dto.ExpCK, dto.ChnlBS)
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

func insertRec(table string) string {
	switch table {
	case poolStructVars:
		return insertStructVar
	case poolLinearVars:
		return insertLinearVar
	default:
		panic("unexpected component-variable table")
	}
}

const (
	poolStructVars = "pool_struct_vars"
	poolLinearVars = "pool_linear_vars"

	insertStructVar = `INSERT INTO pool_struct_vars (comp_id, comp_rn, conn_id, chnl_id, chnl_ph, exp_ck, side) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	insertLinearVar = `INSERT INTO pool_linear_vars (comp_id, comp_rn, conn_id, chnl_id, chnl_ph, exp_ck, side) VALUES ($1, $2, $3, $4, $5, $6, $7)`
)
