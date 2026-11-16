package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/proc/compconn/core"
)

// Package registers the process connection repository.
var Package = do.Package(
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoPgx(logger), nil
	}),
)
