package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/adt/progroot"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	progrootent "orglang/go-engine/internal/pool/ent/progroot"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

// NewDaoEnt creates a pool Program root repository backed by Ent.
func NewDaoEnt(client *ent.Client, log *slog.Logger) progroot.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

func (dao *daoEnt) AddRef(ctx context.Context, conn tx.Conn, ref progroot.RootRef) error {
	err := dao.client.ProgRoot.Create().
		SetID(ref.ProgID.String()).
		SetProgRn(int64(ref.ProgRN)).
		Exec(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", slog.Any("ref", ref))
		return err
	}
	return nil
}

func (dao *daoEnt) TouchRef(ctx context.Context, conn tx.Conn, ref progroot.RootRef) error {
	updated, err := dao.client.ProgRoot.Update().
		AddProgRn(1).
		Where(
			progrootent.ID(ref.ProgID.String()),
			progrootent.ProgRn(int64(ref.ProgRN)),
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
