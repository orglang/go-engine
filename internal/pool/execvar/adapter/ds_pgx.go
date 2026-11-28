package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	adtexecvar "orglang/go-engine/internal/adt/execvar"
	execvaradapter "orglang/go-engine/internal/adt/execvar/adapter"

	execvar "orglang/go-engine/internal/pool/execvar/core"
)

type daoPgx struct {
	log *slog.Logger
}

func NewDaoPgx(log *slog.Logger) execvar.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{log.With(name)}
}

// for compilation purposes
func newRepo() execvar.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRecs(ctx context.Context, conn tx.Conn, recs []adtexecvar.VarRec) (err error) {
	batch := pgx.Batch{}
	for _, rec := range recs {
		dto := execvaradapter.DataFromVarRec(rec)
		batch.Queue(insertRec(getTableName(rec)), dto.ExecID, dto.ExecRN, dto.CommID, dto.ChnlID, dto.ChnlPH, dto.ExpCK, dto.ChnlBS)
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

func getTableName(rec adtexecvar.VarRec) string {
	switch rec.(type) {
	case adtexecvar.StructRec:
		return poolStructVars
	case adtexecvar.LinearRec:
		return poolLinearVars
	default:
		panic(adtexecvar.ErrUnexpectedRecType(rec))
	}
}

func insertRec(table string) string {
	switch table {
	case poolStructVars:
		return insertStructVar
	case poolLinearVars:
		return insertLinearVar
	default:
		panic("unexpected execution-variable table")
	}
}

const (
	poolStructVars = "pool_struct_vars"
	poolLinearVars = "pool_linear_vars"

	insertStructVar = `INSERT INTO pool_struct_vars (exec_id, exec_rn, comm_id, chnl_id, chnl_ph, expr_ck, side) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	insertLinearVar = `INSERT INTO pool_linear_vars (exec_id, exec_rn, comm_id, chnl_id, chnl_ph, expr_ck, side) VALUES ($1, $2, $3, $4, $5, $6, $7)`
)
