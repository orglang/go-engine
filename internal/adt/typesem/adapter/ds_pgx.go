package adapter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
)

type daoPgx struct {
	typeTable string
	descTable string
	log       *slog.Logger
}

func NewDaoPgx(typeTable, descTable string) func(log *slog.Logger) *daoPgx {
	return func(log *slog.Logger) *daoPgx {
		name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
		return &daoPgx{
			typeTable: strings.TrimSpace(typeTable),
			descTable: strings.TrimSpace(descTable),
			log:       log.With(name),
		}
	}
}

// for compilation purposes
func newRepo() typesem.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) TouchRef(ctx context.Context, conn tx.Conn, ref typesem.SemRef) error {
	dto := DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	const query = `UPDATE %s SET type_rn = type_rn + 1 WHERE type_id = $1 AND type_rn = $2`
	sql := fmt.Sprintf(query, dao.typeTable)
	ct, execErr := conn.Pgx.Exec(ctx, sql, dto.TypeID, dto.TypeRN)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	if ct.RowsAffected() == 0 {
		dao.log.Error("touching failed", refAttr)
		return errConcurrentModification(ref)
	}
	dao.log.Log(ctx, lf.LevelTrace, "touching succeed", refAttr)
	return nil
}

func (dao *daoPgx) GetRefsByQNs(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (_ map[uniqsym.ADT]typesem.SemRef, err error) {
	dao.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("qns", typeQNs))
	if len(typeQNs) == 0 {
		return make(map[uniqsym.ADT]typesem.SemRef), nil
	}
	const query = `SELECT t.type_id, t.type_rn
    FROM %s AS t
    JOIN %s AS d ON d.desc_id = t.type_id
    WHERE d.desc_qn = $1`
	sql := fmt.Sprintf(query, dao.typeTable, dao.descTable)
	batch := pgx.Batch{}
	for _, typeQN := range typeQNs {
		batch.Queue(sql, uniqsym.ConvertToString(typeQN))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]SemRef, len(typeQNs))
	for _, typeQN := range typeQNs {
		qnAttr := slog.Any("qn", typeQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr)
			return nil, readErr
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[SemRef])
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return nil, scanErr
		}
		dtos[typeQN] = dto
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return DataToRefMap(dtos)
}

func errConcurrentModification(got typesem.SemRef) error {
	return fmt.Errorf("concurrent modification: %v", got)
}
