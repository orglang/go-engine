package it

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"orglang/go-engine/internal/app/web/api"
)

func TestIntegrationGeneratedEchoServer(t *testing.T) {
	t.Run("registers route and binds path and query parameters", func(t *testing.T) {
		handler := &apiServerStub{}
		e := echo.New()
		api.RegisterHandlers(e, handler)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/processes/executions/exec-123?rn=42", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
		}
		if !handler.called {
			t.Fatal("generated route did not invoke the API handler")
		}
		if handler.id != "exec-123" {
			t.Errorf("unexpected path parameter: got %q, want %q", handler.id, "exec-123")
		}
		if handler.rn == nil || *handler.rn != 42 {
			t.Errorf("unexpected revision query parameter: got %v, want 42", handler.rn)
		}

		var body struct {
			ID string `json:"id"`
			RN int64  `json:"rn"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if body.ID != "exec-123" || body.RN != 42 {
			t.Errorf("unexpected response body: got %+v, want id=exec-123 rn=42", body)
		}
	})

	t.Run("rejects an invalid query parameter before invoking the handler", func(t *testing.T) {
		handler := &apiServerStub{}
		e := echo.New()
		api.RegisterHandlers(e, handler)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/processes/executions/exec-123?rn=not-an-integer", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
		}
		if handler.called {
			t.Fatal("API handler was invoked despite an invalid query parameter")
		}
	})
}

type apiServerStub struct {
	api.ServerInterface

	called bool
	id     string
	rn     *int64
}

func (s *apiServerStub) GetProcessExecution(c echo.Context, id string, params api.GetProcessExecutionParams) error {
	s.called = true
	s.id = id
	s.rn = params.Rn

	if params.Rn == nil {
		return c.JSON(http.StatusOK, map[string]string{"id": id})
	}

	return c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
		RN int64  `json:"rn"`
	}{
		ID: id,
		RN: *params.Rn,
	})
}
