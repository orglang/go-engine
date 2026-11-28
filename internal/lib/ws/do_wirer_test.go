package ws

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/app"
	"orglang/go-engine/internal/lib/kv"
	"orglang/go-engine/internal/lib/lf"
)

type staticExchangeLoader struct{}

func (staticExchangeLoader) Load(key string, dst any) error {
	if key != "exchange" {
		return nil
	}
	config, ok := dst.(*exchangeCS)
	if !ok {
		return errors.New("exchange destination must be *exchangeCS")
	}
	*config = exchangeCS{
		Protocol: protocolCS{
			Modes: []protoModeCS{httpProto},
			HTTP:  httpCS{Port: 8080},
		},
		Server: serverCS{Mode: echoServer},
	}
	return nil
}

func TestPackageProvidesEchoServerAndStopsBeforeStart(t *testing.T) {
	application := app.NewApp(
		lf.Package,
		do.Package(func(i do.Injector) {
			do.ProvideValue[kv.Loader](i, staticExchangeLoader{})
		}),
		Package,
	)

	if _, err := do.Invoke[*echo.Echo](application.Container()); err != nil {
		t.Fatalf("resolve Echo server: %v", err)
	}
	if err := application.Stop(context.Background()); err != nil {
		t.Fatalf("stop application before start: %v", err)
	}
}

func TestNewEchoLogsSuccessfulRequests(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	server := newEcho(logger)
	server.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("health status = %d, want %d", response.Code, http.StatusNoContent)
	}
	output := logs.String()
	for _, want := range []string{"request completed", "method=GET", "uri=/health", "status=204", "latency="} {
		if !strings.Contains(output, want) {
			t.Errorf("request log %q does not contain %q", output, want)
		}
	}
}

func TestEchoServerResourceStartsAndShutsDown(t *testing.T) {
	server := echo.New()
	resource := &echoServerResource{
		server:       server,
		address:      "127.0.0.1:0",
		shutdownDone: make(chan struct{}),
	}
	server.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	if err := resource.start(context.Background()); err != nil {
		t.Fatalf("start server: %v", err)
	}
	address := resource.listenerAddr.String()
	client := &http.Client{Timeout: time.Second}
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://"+address+"/health", nil)
	if err != nil {
		t.Fatalf("create health request: %v", err)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("request health endpoint: %v", err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("health status = %d, want %d", response.StatusCode, http.StatusNoContent)
	}

	if err := resource.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown server: %v", err)
	}
}
