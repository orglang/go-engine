package web

import (
	"context"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/te"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc"
)

type testTransactor struct{}

func (testTransactor) ExplicitTx(context.Context, func(context.Context, tx.Conn) error) error {
	return nil
}

func (testTransactor) ImplicitTx(context.Context, func(context.Context, tx.Conn) error) error {
	return nil
}

func TestPackageProvidesHTTPAdapters(t *testing.T) {
	injector := do.New(
		lf.Package,
		do.Package(func(i do.Injector) {
			do.ProvideValue[tx.Transactor](i, testTransactor{})
		}),
		proc.Package,
		Package,
	)

	if _, err := do.Invoke[te.Renderer](injector); err != nil {
		t.Fatalf("resolve web renderer: %v", err)
	}
	if _, err := do.Invoke[Controller](injector); err != nil {
		t.Fatalf("resolve web controller: %v", err)
	}
}

type testController struct{}

func (testController) Home(*echo.Context) error { return nil }

func TestRegisterRoutes(t *testing.T) {
	server := echo.New()
	RegisterRoutes(server, testController{})

	routes := server.Router().Routes()
	if len(routes) != 1 {
		t.Fatalf("expected one route, got %d", len(routes))
	}
	if routes[0].Method != "GET" || routes[0].Path != "/" {
		t.Fatalf("unexpected route: %s %s", routes[0].Method, routes[0].Path)
	}
}
