package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	compvar "orglang/go-engine/internal/pool/compvar/core"
)

// Package registers the pool component-variable repository.
var Package = do.Package(
	do.Lazy[compvar.Repo](func(i do.Injector) (compvar.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoPgx(logger), nil
	}),
)
