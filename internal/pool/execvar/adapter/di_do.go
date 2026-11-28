package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	execvar "orglang/go-engine/internal/pool/execvar/core"
)

// Package registers the pool execution-variable repository.
var Package = do.Package(
	do.Lazy[execvar.Repo](func(i do.Injector) (execvar.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoPgx(logger), nil
	}),
)
