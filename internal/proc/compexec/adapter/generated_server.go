package adapter

import (
	"github.com/labstack/echo/v4"

	webapi "orglang/go-engine/internal/app/web/api"
)

const processExecutionRoute = "/api/v1/processes/executions"

type generatedRouteFilter struct {
	webapi.EchoRouter
}

func (r generatedRouteFilter) CONNECT(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r generatedRouteFilter) DELETE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r generatedRouteFilter) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	if path != processExecutionRoute+"/:id" {
		return nil
	}
	return r.EchoRouter.GET(path, h, m...)
}
func (r generatedRouteFilter) HEAD(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r generatedRouteFilter) OPTIONS(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r generatedRouteFilter) PATCH(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r generatedRouteFilter) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	if path != processExecutionRoute+"/:id/steps" {
		return nil
	}
	return r.EchoRouter.POST(path, h, m...)
}
func (r generatedRouteFilter) PUT(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r generatedRouteFilter) TRACE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}

// RegisterGeneratedRoutes mounts the migrated generated process-execution routes.
func RegisterGeneratedRoutes(e *echo.Echo, h GeneratedController) error {
	registerGeneratedServer(e, h)
	return nil
}

func cfgGeneratedServer(e *echo.Echo, h *controllerEcho) {
	registerGeneratedServer(e, h)
}

func registerGeneratedServer(e *echo.Echo, h webapi.ServerInterface) {
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, h)
}
