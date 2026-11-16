package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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

type generatedServerAdapter struct {
	webapi.ServerInterface
	handler interface {
		CreateProcessDeclaration(echo.Context) error
	}
}

func (a generatedServerAdapter) CreateProcessDeclaration(c echo.Context) error {
	return a.handler.CreateProcessDeclaration(c)
}

// RegisterGeneratedRoutes mounts the migrated contract-generated process
// declaration route without exposing the other, not-yet-migrated operations.
func RegisterGeneratedRoutes(e *echo.Echo, h Controller) error {
	registerGeneratedServer(e, generatedServerAdapter{handler: h})
	return nil
}

// registerGeneratedServer registers the generated handlers through the route
// filter so only the explicitly migrated operation is exposed.
func registerGeneratedServer(e *echo.Echo, h webapi.ServerInterface) {
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, h)
}

type generatedServerStub struct {
	webapi.ServerInterface
}

func (generatedServerStub) CreateProcessDeclaration(c echo.Context) error {
	return c.NoContent(http.StatusCreated)
}

func TestCfgGeneratedServerRegistersGeneratedRoutes(t *testing.T) {
	e := echo.New()
	registerGeneratedServer(e, generatedServerStub{})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, processDeclarationRoute, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("unexpected status: got %d, want %d", rec.Code, http.StatusCreated)
	}

	req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/programs", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unmigrated operation was registered: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}
