package db

import (
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/app"
	"orglang/go-engine/internal/lib/db/adapter"
	"orglang/go-engine/internal/lib/kv"
	"orglang/go-engine/internal/lib/tx"
)

// Package registers the PostgreSQL transactor and its configuration.
var Package = do.Package(
	do.Lazy[storageProps](func(i do.Injector) (storageProps, error) {
		loader, err := do.Invoke[kv.Loader](i)
		if err != nil {
			return storageProps{}, err
		}
		return newStorageProps(loader)
	}),
	do.Lazy[*pgxPoolResource](func(i do.Injector) (*pgxPoolResource, error) {
		props, err := do.Invoke[storageProps](i)
		if err != nil {
			return nil, err
		}
		lifecycle, err := do.Invoke[app.Lifecycle](i)
		if err != nil {
			return nil, err
		}
		return newPoolPgxDo(props, lifecycle)
	}),
	do.Lazy[tx.Transactor](func(i do.Injector) (tx.Transactor, error) {
		resource, err := do.Invoke[*pgxPoolResource](i)
		if err != nil {
			return nil, err
		}
		return adapter.NewTransactorPgx(resource.pool), nil
	}),
)
