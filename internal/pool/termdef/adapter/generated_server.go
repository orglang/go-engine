package adapter

import (
	"github.com/labstack/echo/v4"

	webapi "orglang/go-engine/internal/app/web/api"
)

const poolDeclarationsRoute = "/api/v1/pools/declarations"

type generatedRouteFilter struct {
	webapi.EchoRouter
}

func (r generatedRouteFilter) CONNECT(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r generatedRouteFilter) DELETE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r generatedRouteFilter) GET(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
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
	if path != poolDeclarationsRoute {
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

func RegisterGeneratedRoutes(e *echo.Echo, h GeneratedController) error {
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, generatedServer{handler: h})
	return nil
}

type generatedServer struct {
	webapi.ServerInterface
	handler GeneratedController
}

func (s generatedServer) CreatePoolDeclaration(c echo.Context) error {
	return s.handler.CreatePoolDeclaration(c)
}
