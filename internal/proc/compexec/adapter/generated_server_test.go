package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	webapi "orglang/go-engine/internal/app/web/api"
)

type generatedServerStub struct {
	webapi.ServerInterface
	gotID string
	gotRN *int64
	stepID string
}

func (s *generatedServerStub) GetProcessExecution(_ echo.Context, id string, params webapi.GetProcessExecutionParams) error {
	s.gotID = id
	s.gotRN = params.Rn
	return nil
}

func (s *generatedServerStub) StepProcessExecution(_ echo.Context, id string) error {
	s.stepID = id
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

	req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, processExecutionRoute+"/exec-2/steps", strings.NewReader(`{"proc_expr":{"kind":"close","conn_ph":"c"}}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || stub.stepID != "exec-2" {
		t.Fatalf("generated step route not registered: status=%d id=%q body=%s", rec.Code, stub.stepID, rec.Body.String())
	}

	req = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/procs/exec-1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("legacy route is still registered: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

