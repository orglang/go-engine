package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/lib/tx"
)

// DAO shares Description root persistence behavior across Ent clients.
type DAO[T any] struct {
	client T
	create func(T, context.Context, tx.Conn, descroot.SemRef) error
	update func(T, context.Context, tx.Conn, descroot.SemRef) (int, error)
	log    *slog.Logger
}

// NewDaoEnt creates a Description root repository using the supplied Ent operations.
func NewDaoEnt[T any](
	client T,
	log *slog.Logger,
	create func(T, context.Context, tx.Conn, descroot.SemRef) error,
	update func(T, context.Context, tx.Conn, descroot.SemRef) (int, error),
) *DAO[T] {
	name := slog.String("name", reflect.TypeFor[DAO[T]]().Name())
	return &DAO[T]{
		client: client,
		create: create,
		update: update,
		log:    log.With(name),
	}
}

// AddRef inserts a new aggregate root at its initial revision.
func (dao *DAO[T]) AddRef(ctx context.Context, conn tx.Conn, ref descroot.SemRef) error {
	if err := dao.create(dao.client, ctx, conn, ref); err != nil {
		dao.log.Error("query execution failed", slog.Any("ref", ref))
		return err
	}
	return nil
}

// TouchRef advances the Description revision if the expected revision matches.
func (dao *DAO[T]) TouchRef(ctx context.Context, conn tx.Conn, ref descroot.SemRef) error {
	updated, err := dao.update(dao.client, ctx, conn, ref)
	if err != nil {
		dao.log.Error("query execution failed", slog.Any("ref", ref))
		return err
	}
	if updated == 0 {
		return fmt.Errorf("concurrent modification: %v", ref)
	}
	return nil
}
