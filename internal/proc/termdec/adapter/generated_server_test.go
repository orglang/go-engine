package adapter

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/progroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

type generatedDeclarationAPIStub struct{}

func (generatedDeclarationAPIStub) Incept(context.Context, uniqsym.ADT) (identity.ADT, progroot.RootRef, error) {
	return identity.ADT{}, progroot.RootRef{}, nil
}
func (generatedDeclarationAPIStub) Create(context.Context, termdec.DecSpec) (termdec.DecSnap, error) {
	return termdec.DecSnap{}, nil
}
func (generatedDeclarationAPIStub) RetrieveSnap(context.Context, identity.ADT) (termdec.DecSnap, error) {
	return termdec.DecSnap{}, nil
}
func (generatedDeclarationAPIStub) RetreiveRefs(context.Context) ([]identity.ADT, error) {
	return nil, nil
}

func TestRegisterGeneratedRoutesRegistersOnlyDeclarationRoutes(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	controller := NewControllerEcho(generatedDeclarationAPIStub{}, logger)
	if err := RegisterGeneratedRoutes(e, controller); err != nil {
		t.Fatalf("RegisterGeneratedRoutes() error = %v", err)
	}

	for _, route := range []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/api/v1/processes/term-1/definition"},
		{method: http.MethodPost, path: "/api/v1/programs"},
	} {
		req := httptest.NewRequestWithContext(context.Background(), route.method, route.path, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("non-Declaration route %s %s status = %d, want %d", route.method, route.path, rec.Code, http.StatusNotFound)
		}
	}
}
