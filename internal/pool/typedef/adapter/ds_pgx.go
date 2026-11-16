package adapter

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/uniqsym"

	typedef "orglang/go-engine/internal/pool/typedef/core"
)

func (dao *daoEnt) GetRecsByQNs(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (_ map[uniqsym.ADT]typedef.DefRec, err error) {
	if len(typeQNs) == 0 {
		return map[uniqsym.ADT]typedef.DefRec{}, nil
	}
	batch := pgx.Batch{}
	for _, typeQN := range typeQNs {
		batch.Queue(selectRecByQN, uniqsym.ConvertToString(typeQN))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]typedef.DefRecDS, len(typeQNs))
	for _, typeQN := range typeQNs {
		qnAttr := slog.Any("qn", typeQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr, slog.Any("sql", selectRecByQN))
			return map[uniqsym.ADT]typedef.DefRec{}, readErr
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[typedef.DefRecDS])
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return map[uniqsym.ADT]typedef.DefRec{}, scanErr
		}
		dtos[typeQN] = dto
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return DataToDefRecMap(dtos)
}

const selectRecByQN = `
	select
		xd.type_id,
		xd.type_rn,
		xd.exp_ck
	from pool_type_defs xd
	join pool_desc_binds db
		on db.desc_id = xd.type_id
	where db.desc_qn = $1`
