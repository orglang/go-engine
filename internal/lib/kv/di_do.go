package kv

import (
	"log/slog"

	"github.com/samber/do/v2"
)

// Package registers shared configuration providers.
var Package = do.Package(
	do.Lazy[Loader](func(i do.Injector) (Loader, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return newViperDriver(logger)
	}),
)
