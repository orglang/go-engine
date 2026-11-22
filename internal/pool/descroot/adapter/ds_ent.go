package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	descrootent "orglang/go-engine/internal/pool/ent/descroot"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

// NewDaoEnt creates a pool Description root repository backed by Ent.
func NewDaoEnt(client *ent.Client, log *slog.Logger) descroot.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

func (dao *daoEnt) AddRef(ctx context.Context, conn tx.Conn, ref descroot.SemRef) error {
	err := dao.client.DescRoot.Create().
		SetID(ref.DescID.String()).
		SetDescRn(int64(ref.DescRN)).
		Exec(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", slog.Any("ref", ref))
		return err
	}
	return nil
}

func (dao *daoEnt) TouchRef(ctx context.Context, conn tx.Conn, ref descroot.SemRef) error {
	updated, err := dao.client.DescRoot.Update().
		AddDescRn(1).
		Where(
			descrootent.ID(ref.DescID.String()),
			descrootent.DescRn(int64(ref.DescRN)),
		).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", slog.Any("ref", ref))
		return err
	}
	if updated == 0 {
		return fmt.Errorf("concurrent modification: %v", ref)
	}
	return nil
}
