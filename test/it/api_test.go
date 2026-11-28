package it

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"

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

	t.Run("binds a JSON request and writes a typed creation response", func(t *testing.T) {
		handler := &apiServerBodyStub{}
		e := echo.New()
		api.RegisterHandlers(e, handler)

		const requestBody = `{
			"asset_vars": [{"comm_ph": "asset", "type_qn": "example.Asset"}],
			"liab_var": {"comm_ph": "liability", "type_qn": "example.Liability"},
			"term_qn": "example.Process"
		}`
		req := httptest.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/api/v1/processes/declarations",
			strings.NewReader(requestBody),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusCreated, rec.Body.String())
		}
		if !handler.called {
			t.Fatal("generated route did not invoke the API handler")
		}
		if handler.body.TermQn != "example.Process" {
			t.Errorf("unexpected term_qn: got %q, want %q", handler.body.TermQn, "example.Process")
		}
		if len(handler.body.AssetVars) != 1 || handler.body.AssetVars[0].CommPh != "asset" {
			t.Errorf("unexpected asset_vars: got %+v", handler.body.AssetVars)
		}
		if handler.body.LiabVar.CommPh != "liability" {
			t.Errorf("unexpected liab_var: got %+v", handler.body.LiabVar)
		}

		var response api.AnyTermRef
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("decode typed response: %v", err)
		}
		if response.TermId != "term-123" || response.ProgRn != 7 {
			t.Errorf("unexpected typed response: got %+v, want term-123 revision 7", response)
		}
	})

	t.Run("maps handler errors to HTTP error responses", func(t *testing.T) {
		e := echo.New()
		api.RegisterHandlers(e, &apiServerErrorStub{})

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/processes/executions/exec-123", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusUnprocessableEntity, rec.Body.String())
		}

		var response struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if response.Message != "execution cannot be advanced" {
			t.Errorf("unexpected error message: got %q", response.Message)
		}
	})

	t.Run("supports a no-content operation", func(t *testing.T) {
		e := echo.New()
		api.RegisterHandlers(e, &apiServerNoContentStub{})

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/pools/executions/steps", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusNoContent, rec.Body.String())
		}
		if rec.Body.Len() != 0 {
			t.Errorf("expected an empty response body, got %q", rec.Body.String())
		}
	})
}

type apiServerStub struct {
	api.ServerInterface

	called bool
	id     string
	rn     *int64
}

func (s *apiServerStub) GetProcessExecution(c *echo.Context, id string, params api.GetProcessExecutionParams) error {
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

type apiServerBodyStub struct {
	api.ServerInterface

	called bool
	body   api.ProcDecSpec
}

func (s *apiServerBodyStub) CreateProcessDeclaration(c *echo.Context) error {
	s.called = true
	if err := c.Bind(&s.body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid declaration")
	}

	return c.JSON(http.StatusCreated, api.AnyTermRef{
		TermId: "term-123",
		ProgRn: 7,
	})
}

type apiServerErrorStub struct {
	api.ServerInterface
}

func (s *apiServerErrorStub) GetProcessExecution(*echo.Context, string, api.GetProcessExecutionParams) error {
	return echo.NewHTTPError(http.StatusUnprocessableEntity, "execution cannot be advanced")
}

type apiServerNoContentStub struct {
	api.ServerInterface
}

func (s *apiServerNoContentStub) StepPoolExecution(c *echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
