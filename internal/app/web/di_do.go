package web

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/lib/te"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

// Package registers the web application's renderer and HTTP controller.
var Package = do.Package(
	do.Lazy[te.Renderer](func(i do.Injector) (te.Renderer, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return newRendererStdlib(logger)
	}),
	do.Lazy[Controller](func(i do.Injector) (Controller, error) {
		api, err := do.Invoke[typedef.API](i)
		if err != nil {
			return nil, err
		}
		renderer, err := do.Invoke[te.Renderer](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return newControllerEcho(api, renderer, logger), nil
	}),
)
