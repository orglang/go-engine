package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/proc/typeexp/core"
)

// Package registers process type-expression persistence.
var Package = do.Package(
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoPgx(logger), nil
	}),
)
