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

type execRootDaoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewExecRootDaoEnt(client *ent.Client, log *slog.Logger) execroot.Repo {
	name := slog.String("name", reflect.TypeFor[execRootDaoEnt]().Name())
	return &execRootDaoEnt{client: client, log: log.With(name)}
}

func (dao *execRootDaoEnt) TouchRef(ctx context.Context, conn tx.Conn, ref execroot.RootRef) error {
	dto := execrootadapter.DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	updated, err := dao.client.ExecEval.Update().
		AddExecRn(1).
		Where(execeval.ID(dto.ExecID), execeval.ExecRn(dto.ExecRN)).
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
