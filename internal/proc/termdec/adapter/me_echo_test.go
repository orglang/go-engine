package adapter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/uniqsym"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

func TestCreateProcessDeclarationHandler(t *testing.T) {
	t.Run("creates a declaration and returns its identity", func(t *testing.T) {
		decAPI := &generatedDecAPIStub{}
		e := echo.New()
		h := &controllerEcho{api: decAPI}
		e.POST("/api/v1/processes/declarations", h.CreateProcessDeclaration)

		const body = `{
			"asset_vars": [{"conn_ph": "asset", "type_qn": "example.Asset"}],
			"liab_var": {"conn_ph": "liability", "type_qn": "example.Liability"},
			"term_qn": "example.Process"
		}`
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost,
			"/api/v1/processes/declarations", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusCreated, rec.Body.String())
		}
		if !decAPI.called {
			t.Fatal("declaration API was not called")
		}
		if got := uniqsym.ConvertToString(decAPI.spec.TermQN); got != "example.Process" {
			t.Errorf("unexpected term qn: got %q", got)
		}
		if got := uniqsym.ConvertToString(decAPI.spec.LiabVar.TypeQN); got != "example.Liability" {
			t.Errorf("unexpected liability type qn: got %q", got)
		}
		if got := string(decAPI.spec.LiabVar.ChnlPH); got != "liability" {
			t.Errorf("unexpected liability channel: got %q", got)
		}
		if len(decAPI.spec.AssetVars) != 1 || uniqsym.ConvertToString(decAPI.spec.AssetVars[0].TypeQN) != "example.Asset" {
			t.Errorf("unexpected asset vars: %+v", decAPI.spec.AssetVars)
		}

		var response struct {
			TermID string `json:"term_id"`
			TermRn int64  `json:"term_rn"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if response.TermRn != 11 {
			t.Errorf("unexpected term revision: got %d, want 11", response.TermRn)
		}
	})

	t.Run("rejects invalid qualified names before creating", func(t *testing.T) {
		decAPI := &generatedDecAPIStub{}
		e := echo.New()
		h := &controllerEcho{api: decAPI}
		e.POST("/api/v1/processes/declarations", h.CreateProcessDeclaration)

		const body = `{
			"asset_vars": [],
			"liab_var": {"conn_ph": "liability", "type_qn": ""},
			"term_qn": "example.Process"
		}`
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost,
			"/api/v1/processes/declarations", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("unexpected status: got %d, want %d; body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
		}
		if decAPI.called {
			t.Fatal("declaration API was called for an invalid request")
		}
	})
}

type generatedDecAPIStub struct {
	termdec.API
	called bool
	spec   termdec.DecSpec
}

func (s *generatedDecAPIStub) Create(spec termdec.DecSpec) (termdec.DecSnap, error) {
	s.called = true
	s.spec = spec
	return termdec.DecSnap{TermRef: termsem.SemRef{TermRN: 11}}, nil
}
