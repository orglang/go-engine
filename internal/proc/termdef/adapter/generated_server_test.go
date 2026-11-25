package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"

	webapi "orglang/go-engine/internal/app/web/api"
)

type generatedServerStub struct {
	webapi.ServerInterface
	gotID           string
	gotRN           *int64
	gotDefinitionID string
	gotDefinitionRN *int64
}

func (s *generatedServerStub) CreateProcessDeclaration(c *echo.Context) error {
	return c.NoContent(http.StatusCreated)
}

func (s *generatedServerStub) GetProcessDefinition(c *echo.Context, id string, params webapi.GetProcessDefinitionParams) error {
	s.gotDefinitionID = id
	s.gotDefinitionRN = params.Rn
	return c.NoContent(http.StatusOK)
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

	req = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/processes/proc-2/definition?rn=9", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected Definition GET status: got %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if stub.gotDefinitionID != "proc-2" {
		t.Errorf("unexpected Definition id: got %q, want %q", stub.gotDefinitionID, "proc-2")
	}
	if stub.gotDefinitionRN == nil || *stub.gotDefinitionRN != 9 {
		t.Errorf("unexpected Definition rn: got %v, want 9", stub.gotDefinitionRN)
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
