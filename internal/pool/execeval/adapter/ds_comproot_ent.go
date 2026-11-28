package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/adt/execroot"
	execrootadapter "orglang/go-engine/internal/adt/execroot/adapter"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	"orglang/go-engine/internal/pool/ent/execeval"
)

type compRootDaoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewCompRootDaoEnt(client *ent.Client, log *slog.Logger) execroot.Repo {
	name := slog.String("name", reflect.TypeFor[compRootDaoEnt]().Name())
	return &compRootDaoEnt{client: client, log: log.With(name)}
}

func (dao *compRootDaoEnt) TouchRef(ctx context.Context, conn tx.Conn, ref execroot.RootRef) error {
	dto := execrootadapter.DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	updated, err := dao.client.ExecEval.Update().
		AddCompRn(1).
		Where(execeval.ID(dto.RootID), execeval.CompRn(dto.RootRN)).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}
	if updated == 0 {
		dao.log.Error("touching failed", refAttr)
		return fmt.Errorf("concurrent modification: %v", ref)
	}
	dao.log.Log(ctx, lf.LevelTrace, "touching succeed", refAttr)
	return nil
}
