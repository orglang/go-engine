package adapter

import (
	"github.com/labstack/echo/v5"

	webapi "orglang/go-engine/internal/app/web/api"
)

const processTypeDefinitionsRoute = "/api/v1/processes/types"

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
	if path != processTypeDefinitionsRoute+"/:id" {
		return nil
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
	if path != processTypeDefinitionsRoute {
		return nil
	}
	return r.EchoRouter.POST(path, h, m...)
}

func (r generatedRouteFilter) PUT(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}

func (r generatedRouteFilter) TRACE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}

func RegisterGeneratedRoutes(e *echo.Echo, h GeneratedController) error {
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, generatedServer{handler: h})
	return nil
}

type generatedServer struct {
	webapi.ServerInterface
	handler GeneratedController
}

func (s generatedServer) CreateProcessTypeDefinition(c *echo.Context) error {
	return s.handler.CreateProcessTypeDefinition(c)
}

func (s generatedServer) GetProcessTypeDefinition(c *echo.Context, id string, params webapi.GetProcessTypeDefinitionParams) error {
	return s.handler.GetProcessTypeDefinition(c, id, params)
}
