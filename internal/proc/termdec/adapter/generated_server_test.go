package adapter

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/uniqsym"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

type generatedDeclarationAPIStub struct{}

func (generatedDeclarationAPIStub) Incept(context.Context, uniqsym.ADT) (termsem.SemRef, error) {
	return termsem.SemRef{}, nil
}
func (generatedDeclarationAPIStub) Create(context.Context, termdec.DecSpec) (termdec.DecSnap, error) {
	return termdec.DecSnap{}, nil
}
func (generatedDeclarationAPIStub) RetrieveSnap(context.Context, termsem.SemRef) (termdec.DecSnap, error) {
	return termdec.DecSnap{}, nil
}
func (generatedDeclarationAPIStub) RetreiveRefs(context.Context) ([]termsem.SemRef, error) {
	return nil, nil
}

func TestRegisterGeneratedRoutesRegistersOnlyDeclarationRoutes(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	controller := NewControllerEcho(generatedDeclarationAPIStub{}, logger)
	if err := RegisterGeneratedRoutes(e, controller); err != nil {
		t.Fatalf("RegisterGeneratedRoutes() error = %v", err)
	}

	routes := map[string]bool{}
	for _, route := range e.Routes() {
		routes[route.Method+" "+route.Path] = true
	}
	for _, route := range []string{
		"POST " + processDeclarationRoute,
		"GET " + processDeclarationRoute + "/:id",
	} {
		if !routes[route] {
			t.Errorf("Declaration route %q was not registered", route)
		}
	}
	for _, route := range []string{
		"GET /api/v1/processes/:id/definition",
		"POST /api/v1/programs",
	} {
		if routes[route] {
			t.Errorf("non-Declaration route %q was registered", route)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/processes/term-1/definition", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("Definition route status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
