package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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

type generatedServerStub struct {
	webapi.ServerInterface
	gotID string
	gotRN *int64
}

func (s *generatedServerStub) CreateProcessDeclaration(c *echo.Context) error {
	return c.NoContent(http.StatusCreated)
}

func (s *generatedServerStub) GetProcessDeclaration(c *echo.Context, id string, params webapi.GetProcessDeclarationParams) error {
	s.gotID = id
	s.gotRN = params.Rn
	return c.NoContent(http.StatusOK)
}

func TestCfgGeneratedServerRegistersGeneratedRoutes(t *testing.T) {
	e := echo.New()
	stub := &generatedServerStub{}
	registerGeneratedServer(e, stub)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, processDeclarationRoute, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("unexpected POST status: got %d, want %d", rec.Code, http.StatusCreated)
	}

	req = httptest.NewRequestWithContext(context.Background(), http.MethodGet, processDeclarationRoute+"/proc-1?rn=7", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected GET status: got %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if stub.gotID != "proc-1" {
		t.Errorf("unexpected id: got %q, want %q", stub.gotID, "proc-1")
	}
	if stub.gotRN == nil || *stub.gotRN != 7 {
		t.Errorf("unexpected rn: got %v, want 7", stub.gotRN)
	}

	req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/programs", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unmigrated operation was registered: got %d, want %d", rec.Code, http.StatusNotFound)
	}

	for _, route := range []struct {
		method string
		path   string
	}{
		{method: http.MethodPost, path: "/api/v1/procs/decs"},
		{method: http.MethodGet, path: "/api/v1/procs/decs/proc-1"},
	} {
		req = httptest.NewRequestWithContext(context.Background(), route.method, route.path, nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("legacy route %s %s is still registered: got %d, want %d", route.method, route.path, rec.Code, http.StatusNotFound)
		}
	}
}
