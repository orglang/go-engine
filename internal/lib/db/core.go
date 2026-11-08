package db

import (
	"fmt"

	"go.uber.org/fx"

	"orglang/go-engine/lib/db/adapter"
	"orglang/go-engine/lib/tx"
)

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
