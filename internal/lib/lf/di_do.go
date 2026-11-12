package lf

import (
	"log/slog"

	"github.com/samber/do/v2"
)

// Package registers shared logging providers.
var Package = do.Package(
	do.Lazy[*slog.Logger](func(do.Injector) (*slog.Logger, error) {
		return newSlogLogger(), nil
	}),
)
