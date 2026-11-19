package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"

	webapi "orglang/go-engine/internal/app/web/api"
)

const programRoute = "/api/v1/programs"

// generatedRouteFilter registers only the generated operations migrated so far.
type generatedRouteFilter struct {
	webapi.EchoRouter
}

func (r generatedRouteFilter) CONNECT(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}
func (r generatedRouteFilter) DELETE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
}
func (r generatedRouteFilter) GET(string, echo.HandlerFunc, ...echo.MiddlewareFunc) echo.RouteInfo {
	return echo.RouteInfo{}
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
func (r generatedRouteFilter) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	if path != programRoute {
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

func cfgGeneratedServer(e *echo.Echo, h *controllerEcho) {
	registerGeneratedServer(e, h)
}

func registerGeneratedServer(e *echo.Echo, h webapi.ServerInterface) {
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, h)
}

type generatedServerStub struct {
	webapi.ServerInterface
	called bool
}

func (s *generatedServerStub) CreateProgram(c *echo.Context) error {
	s.called = true
	return c.NoContent(http.StatusCreated)
}

func TestCfgGeneratedServerRegistersGeneratedProgramRoute(t *testing.T) {
	e := echo.New()
	stub := &generatedServerStub{}
	registerGeneratedServer(e, stub)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, programRoute, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if !stub.called {
		t.Fatal("generated CreateProgram handler was not called")
	}

	req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/progs", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("legacy route is still registered: got %d, want %d", rec.Code, http.StatusNotFound)
	}

	req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/processes/declarations", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unmigrated route was registered: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}
