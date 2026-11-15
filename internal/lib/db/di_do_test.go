package db

import (
	"context"
	"testing"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/app"
	"orglang/go-engine/internal/lib/kv"
)

type staticStorageLoader struct{}

func (staticStorageLoader) Load(key string, dst any) error {
	if key != "storage" {
		return nil
	}
	props := dst.(*storageProps)
	*props = storageProps{
		Protocol: protoProps{
			Mode: postgresMode,
			Postgres: postgresProps{
				URL: "postgres://user:pass@127.0.0.1:5432/engine",
			},
		},
		Driver: driverProps{
			Mode: pgxMode,
			Pgx:  pgxProps{MaxConns: 1},
		},
	}
	return nil
}

func TestPackageProvidesTransactorAndClosesPoolOnAppStop(t *testing.T) {
	application := app.NewApp(
		do.Package(func(i do.Injector) {
			do.ProvideValue[kv.Loader](i, staticStorageLoader{})
		}),
		Package,
	)

	if _, err := do.Invoke[interface{ ExplicitTx(context.Context, func(context.Context, tx.Conn) error) error }](application.Container()); err == nil {
		t.Fatal("expected unsupported interface binding to be absent")
	}

	resource, err := do.Invoke[*pgxPoolResource](application.Container())
	if err != nil {
		t.Fatalf("resolve pool resource: %v", err)
	}
	if _, err := do.Invoke[tx.Transactor](application.Container()); err != nil {
		t.Fatalf("resolve transactor: %v", err)
	}

	if err := application.Stop(context.Background()); err != nil {
		t.Fatalf("stop application: %v", err)
	}
	if _, err := resource.pool.Acquire(context.Background()); err == nil {
		t.Fatal("expected pool to be closed by application shutdown")
	}
}
