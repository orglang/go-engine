package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/typesem"
	typesemadapter "orglang/go-engine/internal/adt/typesem/adapter"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	)

type typeSemDaoPgx struct {
	log *slog.Logger
}

func NewTypeSemDaoPgx(log *slog.Logger) typesem.Repo {
	name := slog.String("name", reflect.TypeFor[typeSemDaoPgx]().Name())
	return &typeSemDaoPgx{log: log.With(name)}
}


func (dao *typeSemDaoPgx) GetRefsByQNs(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (_ map[uniqsym.ADT]typesem.SemRef, err error) {
	dao.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("qns", typeQNs))
	if len(typeQNs) == 0 {
		return make(map[uniqsym.ADT]typesem.SemRef), nil
	}
	const query = `
SELECT t.type_id, t.desc_id, dr.desc_rn AS desc_rn
FROM proc_type_defs AS t
JOIN proc_desc_roots AS dr ON dr.desc_id = t.desc_id
JOIN proc_desc_binds AS d ON d.desc_id = t.type_id
WHERE d.desc_qn = $1
`
	batch := pgx.Batch{}
	for _, typeQN := range typeQNs {
		batch.Queue(query, uniqsym.ConvertToString(typeQN))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]typesemadapter.SemRef, len(typeQNs))
	for _, typeQN := range typeQNs {
		qnAttr := slog.Any("qn", typeQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr)
			return nil, readErr
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[typesemadapter.SemRef])
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return nil, scanErr
		}
		dtos[typeQN] = dto
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return typesemadapter.DataToRefMap(dtos)
}
