package web

import (
	"github.com/labstack/echo/v4"

	webapi "orglang/go-engine/internal/app/web/api"
)

const processDeclarationRoute = "/api/v1/processes/declarations"

// generatedRouteFilter limits generated registration to the operation migrated
// in this PR. Other contract operations remain available for later migrations.
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
	if path != processDeclarationRoute {
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

// cfgGeneratedServer registers only the contract-generated route migrated by
// this PR on the application's existing Echo instance. It does not start a
// second HTTP listener or expose unrelated, not-yet-migrated operations.
func cfgGeneratedServer(e *echo.Echo, h webapi.ServerInterface) {
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, h)
}
