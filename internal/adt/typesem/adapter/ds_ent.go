package adapter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	"orglang/go-engine/internal/proc/ent/typedef"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) typesem.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() typesem.Repo { return new(daoEnt) }

func (dao *daoEnt) TouchRef(ctx context.Context, conn tx.Conn, ref typesem.SemRef) error {
	dto := DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	affected, err := dao.client.TypeDef.Update().
		AddTypeRn(1).
		Where(typedef.ID(dto.TypeID), typedef.TypeRn(dto.TypeRN)).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}
	if affected == 0 {
		dao.log.Error("touching failed", refAttr)
		return errConcurrentModification(ref)
	}
	dao.log.Log(ctx, lf.LevelTrace, "touching succeed", refAttr)
	return nil
}

func (dao *daoEnt) GetRefsByQNs(ctx context.Context, conn tx.Conn, typeQNs []uniqsym.ADT) (_ map[uniqsym.ADT]typesem.SemRef, err error) {
	dao.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("qns", typeQNs))
	if len(typeQNs) == 0 {
		return make(map[uniqsym.ADT]typesem.SemRef), nil
	}
	const query = `
SELECT t.type_id, t.type_rn
FROM proc_type_defs AS t
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
