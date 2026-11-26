package adapter

import (
	"github.com/labstack/echo/v5"

	webapi "orglang/go-engine/internal/app/web/api"
)

// generatedRouteFilter limits generated registration to the Definition read
// operation. Other contract operations remain available for later migrations.
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
	if path != "/api/v1/processes/:id/definition" {
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
func (generatedRouteFilter) POST(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}
func (r generatedRouteFilter) PUT(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}
func (r generatedRouteFilter) TRACE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}

type generatedServer struct {
	webapi.ServerInterface
	controller *DefinitionController
}

func (h generatedServer) GetProcessDefinition(c *echo.Context, id string, params webapi.GetProcessDefinitionParams) error {
	return h.controller.GetProcessDefinition(c, id, params)
}

// RegisterGeneratedRoutes mounts the generated process Definition-read operation.
func RegisterGeneratedRoutes(e *echo.Echo, h *DefinitionController) error {
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, generatedServer{controller: h})
	return nil
}

// registerGeneratedServer is retained for focused generated-router contract tests.
func registerGeneratedServer(e *echo.Echo, h webapi.ServerInterface) {
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, h)
}
