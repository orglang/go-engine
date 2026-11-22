package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	"orglang/go-engine/internal/pool/ent/compexec"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) *daoEnt {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() compsem.Repo { return new(daoEnt) }

func (dao *daoEnt) TouchRef(ctx context.Context, conn tx.Conn, ref compsem.SemRef) error {
	dto := DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	affected, err := dao.client.CompExec.Update().
		AddCompRn(1).
		Where(compexec.ID(dto.CompID), compexec.CompRn(dto.CompRN)).
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

func errConcurrentModification(got compsem.SemRef) error {
	return fmt.Errorf("concurrent modification: %v", got)
}
