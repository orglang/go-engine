package adapter

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

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
func (r generatedRouteFilter) POST(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r generatedRouteFilter) PUT(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r generatedRouteFilter) TRACE(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route {
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

func (s *generatedServerStub) GetProcessExecution(_ echo.Context, id string, params webapi.GetProcessExecutionParams) error {
	s.gotID = id
	s.gotRN = params.Rn
	return nil
}

func TestCfgGeneratedServerRegistersProcessExecutionRoute(t *testing.T) {
	e := echo.New()
	stub := &generatedServerStub{}
	registerGeneratedServer(e, stub)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, processExecutionRoute+"/exec-1?rn=42", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if stub.gotID != "exec-1" {
		t.Errorf("unexpected id: got %q, want %q", stub.gotID, "exec-1")
	}
	if stub.gotRN == nil || *stub.gotRN != 42 {
		t.Errorf("unexpected rn: got %v, want 42", stub.gotRN)
	}

	req = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/procs/exec-1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("legacy route is still registered: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

type legacyStepStub struct {
	called bool
}

func (s *legacyStepStub) PostStep(c echo.Context) error {
	s.called = true
	return c.NoContent(http.StatusOK)
}

func TestCfgLegacyStepRouteRegistersOnlyUnmigratedStepOperation(t *testing.T) {
	e := echo.New()
	handler := &controllerEcho{log: slog.Default()}
	cfgLegacyStepRoute(e, handler)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/procs/exec-1/steps", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code == http.StatusNotFound {
		t.Fatal("legacy step route was not registered")
	}
}
