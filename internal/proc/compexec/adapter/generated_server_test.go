package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	webapi "orglang/go-engine/internal/app/web/api"
)

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
}

func TestRegisterLegacyProcessRoute(t *testing.T) {
	e := echo.New()
	stub := &generatedServerStub{}
	RegisterLegacyProcessRoute(e, stub)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/procs/exec-1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if stub.gotID != "exec-1" {
		t.Errorf("unexpected id: got %q, want %q", stub.gotID, "exec-1")
	}
	if stub.gotRN != nil {
		t.Errorf("unexpected rn: got %v, want nil", *stub.gotRN)
	}
}

func TestCfgLegacyStepRouteRegistersOnlyUnmigratedStepOperation(t *testing.T) {
	e := echo.New()
	cfgLegacyStepRoute(e, &controllerEcho{})

	for _, route := range e.Routes() {
		if route.Method == http.MethodPost && route.Path == "/api/v1/procs/:id/steps" {
			return
		}
	}
	t.Fatal("legacy step route was not registered")
}
