package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"

	webapi "orglang/go-engine/internal/app/web/api"
)

type generatedControllerStub struct {
	gotID string
	gotRN *int64
}

func (s *generatedControllerStub) GetPoolDefinition(c *echo.Context, id string, params webapi.GetPoolDefinitionParams) error {
	s.gotID = id
	s.gotRN = params.Rn
	return c.NoContent(http.StatusOK)
}

func TestRegisterGeneratedRoutesRegistersDefinitionGET(t *testing.T) {
	e := echo.New()
	stub := &generatedControllerStub{}
	if err := RegisterGeneratedRoutes(e, stub); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/pools/pool-1/definition?rn=4", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if stub.gotID != "pool-1" {
		t.Errorf("unexpected id: got %q, want %q", stub.gotID, "pool-1")
	}
	if stub.gotRN == nil || *stub.gotRN != 4 {
		t.Errorf("unexpected rn: got %v, want 4", stub.gotRN)
	}

	req = httptest.NewRequestWithContext(context.Background(), http.MethodPatch, "/api/v1/pools/pool-1/definition?rn=4", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("PATCH was registered unexpectedly: got %d", rec.Code)
	}
}
