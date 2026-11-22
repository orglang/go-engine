package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/adt/compsem"
	compsemadapter "orglang/go-engine/internal/adt/compsem/adapter"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	"orglang/go-engine/internal/pool/ent/compexec"
)

type compSemDaoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewCompSemDaoEnt(client *ent.Client, log *slog.Logger) compsem.Repo {
	name := slog.String("name", reflect.TypeFor[compSemDaoEnt]().Name())
	return &compSemDaoEnt{client: client, log: log.With(name)}
}

func (dao *compSemDaoEnt) TouchRef(ctx context.Context, conn tx.Conn, ref compsem.SemRef) error {
	dto := compsemadapter.DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	updated, err := dao.client.CompExec.Update().
		AddCompRn(1).
		Where(compexec.ID(dto.CompID), compexec.CompRn(dto.CompRN)).
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
