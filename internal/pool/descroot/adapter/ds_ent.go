package adapter

import (
	"context"
	"log/slog"

	"orglang/go-engine/internal/adt/descroot"
	descrootadapter "orglang/go-engine/internal/adt/descroot/adapter"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	descrootent "orglang/go-engine/internal/pool/ent/descroot"
)

// NewDaoEnt creates a pool Description root repository backed by Ent.
func NewDaoEnt(client *ent.Client, log *slog.Logger) descroot.Repo {
	return descrootadapter.NewDaoEnt(
		client,
		log,
		func(client *ent.Client, ctx context.Context, conn tx.Conn, ref descroot.SemRef) error {
			return client.DescRoot.Create().
				SetID(ref.DescID.String()).
				SetDescRn(int64(ref.DescRN)).
				Exec(entpgx.WithConn(ctx, conn.Pgx))
		},
		func(client *ent.Client, ctx context.Context, conn tx.Conn, ref descroot.SemRef) (int, error) {
			return client.DescRoot.Update().
				AddDescRn(1).
				Where(
					descrootent.ID(ref.DescID.String()),
					descrootent.DescRn(int64(ref.DescRN)),
				).
				Save(entpgx.WithConn(ctx, conn.Pgx))
		},
	)
}
