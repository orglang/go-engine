package db

import (
	"context"
	"fmt"

	"orglang/go-engine/lib/db/adapter"
	"orglang/go-engine/lib/tx"

	"go.uber.org/fx"
)

// UoW is kept as an alias for source compatibility.
//
// Deprecated: use tx.Conn and pass context separately.
type UoW = tx.UoW

// Transactor is the public transaction contract.
//
// Deprecated: use tx.Transactor.
type Transactor = tx.Transactor

func newTransactor(props storageProps, lc fx.Lifecycle) (tx.Transactor, error) {
	switch {
	case props.Protocol.Mode == postgresMode && props.Driver.Mode == pgxMode:
		pool, err := newPoolPgx(props, lc)
		if err != nil {
			return nil, err
		}
		return adapter.NewTransactorPgx(pool), nil
	default:
		panic(fmt.Sprintf("unsupported combination: %s & %s", props.Protocol.Mode, props.Driver.Mode))
	}
}
