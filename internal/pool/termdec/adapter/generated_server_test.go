package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

type generatedDeclarationControllerStub struct{}

func (generatedDeclarationControllerStub) CreatePoolDeclaration(c *echo.Context) error {
	return c.NoContent(http.StatusCreated)
}

func TestRegisterGeneratedRoutesRegistersDeclarationPOSTOnly(t *testing.T) {
	e := echo.New()
	if err := RegisterGeneratedRoutes(e, generatedDeclarationControllerStub{}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, poolDeclarationsRoute, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d; body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	req = httptest.NewRequestWithContext(context.Background(), http.MethodGet, poolDeclarationsRoute, nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET route was registered unexpectedly: got %d", rec.Code)
	}
}
