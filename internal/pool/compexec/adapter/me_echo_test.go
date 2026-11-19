package adapter

import (
	"io"
	"log/slog"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestRegisterRoutes(t *testing.T) {
	server := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	controller := newControllerEcho(nil, logger)

	if err := RegisterRoutes(server, controller); err != nil {
		t.Fatalf("register routes: %v", err)
	}

	want := map[string]bool{
		"POST /api/v1/pools/executions":        false,
		"POST /api/v1/pools/executions/steps":  false,
		"POST /api/v1/pools/executions/spawns": false,
	}
	for _, route := range server.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for route, found := range want {
		if !found {
			t.Errorf("route %q was not registered", route)
		}
	}
}
