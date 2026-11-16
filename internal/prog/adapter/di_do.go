package adapter

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"

	webapi "orglang/go-engine/internal/app/web/api"
	coreprog "orglang/go-engine/internal/prog/core"
)

// Controller is the generated API contract implemented by the program adapter.
// It is a distinct port so it can coexist with other generated controllers in
// the same application injector.
type Controller interface {
	webapi.ServerInterface
}

// Package registers the program service and its generated HTTP adapter.
var Package = do.Package(
	do.Lazy[coreprog.API](func(i do.Injector) (coreprog.API, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return coreprog.NewService(logger), nil
	}),
	do.Lazy[Controller](func(i do.Injector) (Controller, error) {
		api, err := do.Invoke[coreprog.API](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return newControllerEcho(api, logger), nil
	}),
)

// RegisterRoutes mounts the generated program endpoint on the application's
// existing Echo server. It does not start an HTTP listener.
func RegisterRoutes(e *echo.Echo, h Controller) error {
	registerGeneratedServer(e, h)
	return nil
}
