package adapter

import (
	"github.com/labstack/echo/v5"

	webapi "orglang/go-engine/internal/app/web/api"
)

const processDeclarationRoute = "/api/v1/processes/declarations"

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
	if path != processDeclarationRoute+"/:id" {
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

type generatedServer struct {
	webapi.ServerInterface
	controller *Controller
}

func (h generatedServer) CreateProcessDeclaration(c *echo.Context) error {
	return h.controller.CreateProcessDeclaration(c)
}

func (h generatedServer) GetProcessDeclaration(c *echo.Context, id string, params webapi.GetProcessDeclarationParams) error {
	return h.controller.GetProcessDeclaration(c, id, params)
}

// RegisterGeneratedRoutes mounts the generated Declaration HTTP operations.
func RegisterGeneratedRoutes(e *echo.Echo, h *Controller) error {
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, generatedServer{controller: h})
	return nil
}
