package adapter

import (
	"github.com/labstack/echo/v5"

	webapi "orglang/go-engine/internal/app/web/api"
)

const processDeclarationRoute = "/api/v1/processes/declarations"

// generatedRouteFilter limits generated registration to the operations migrated
// so far. Other contract operations remain available for later migrations.
type generatedRouteFilter struct {
	webapi.EchoRouter
}

func (r generatedRouteFilter) CONNECT(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}
func (r generatedRouteFilter) DELETE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}
func (r generatedRouteFilter) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	if path != processDeclarationRoute+"/:id" && path != "/api/v1/processes/:id/definition" {
		return echo.RouteInfo{}
	}
	return r.EchoRouter.GET(path, h, m...)
}
func (r generatedRouteFilter) HEAD(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}
func (r generatedRouteFilter) OPTIONS(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}
func (r generatedRouteFilter) PATCH(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}
func (r generatedRouteFilter) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	if path != processDeclarationRoute {
		return echo.RouteInfo{}
	}
	return r.EchoRouter.POST(path, h, m...)
}
func (r generatedRouteFilter) PUT(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}
func (r generatedRouteFilter) TRACE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}

// RegisterGeneratedRoutes mounts the migrated generated process
// declaration and Definition-read operations on the application's existing Echo server.
func RegisterGeneratedRoutes(e *echo.Echo, h webapi.ServerInterface) error {
	registerGeneratedServer(e, h)
	return nil
}

func cfgGeneratedServer(e *echo.Echo, h *controllerEcho) {
	registerGeneratedServer(e, h)
}

func registerGeneratedServer(e *echo.Echo, h webapi.ServerInterface) {
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, h)
}

