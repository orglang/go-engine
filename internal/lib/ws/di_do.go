package ws

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/app"
	"orglang/go-engine/internal/lib/kv"
)

// Package registers the Echo HTTP server and its explicit application lifecycle.
var Package = do.Package(
	do.Lazy[exchangeCS](func(i do.Injector) (exchangeCS, error) {
		loader, err := do.Invoke[kv.Loader](i)
		if err != nil {
			return exchangeCS{}, err
		}
		return newExchangeCS(loader)
	}),
	do.Lazy[*echoServerResource](func(i do.Injector) (*echoServerResource, error) {
		config, err := do.Invoke[exchangeCS](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		lifecycle, err := do.Invoke[app.Lifecycle](i)
		if err != nil {
			return nil, err
		}

		server := newEcho(logger)
		resource := &echoServerResource{
			server:       server,
			address:      fmt.Sprintf(":%d", config.Protocol.HTTP.Port),
			shutdownDone: make(chan struct{}),
		}
		appendErr := lifecycle.Append(app.Hook{
			OnStart: resource.start,
			OnStop:  resource.Shutdown,
		})
		if appendErr != nil {
			return nil, appendErr
		}
		return resource, nil
	}),
	do.Lazy[*echo.Echo](func(i do.Injector) (*echo.Echo, error) {
		resource, err := do.Invoke[*echoServerResource](i)
		if err != nil {
			return nil, err
		}
		return resource.server, nil
	}),
)

func newEcho(logger *slog.Logger) *echo.Echo {
	server := echo.New()
	log := logger.With(slog.String("name", "echoServer"))
	server.Logger = log
	server.Use(middleware.Recover())
	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:   true,
		LogURI:      true,
		LogStatus:   true,
		HandleError: true,
		LogValuesFunc: func(_ *echo.Context, values middleware.RequestLoggerValues) error {
			if values.Error != nil {
				log.Error("request processing failed",
					slog.String("method", values.Method),
					slog.String("uri", values.URI),
					slog.Int("status", values.Status),
					slog.String("reason", values.Error.Error()),
				)
			}
			return nil
		},
	}))
	return server
}

type echoServerResource struct {
	server       *echo.Echo
	address      string
	cancel       context.CancelFunc
	serveDone    chan error
	listenerAddr net.Addr

	shutdownOnce sync.Once
	shutdownDone chan struct{}
	shutdownErr  error
}

func (r *echoServerResource) start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	serveCtx, cancel := context.WithCancel(ctx)
	ready := make(chan struct{})
	r.cancel = cancel
	r.serveDone = make(chan error, 1)

	config := echo.StartConfig{
		Address:         r.address,
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: 10 * time.Second,
		ListenerAddrFunc: func(addr net.Addr) {
			r.listenerAddr = addr
			close(ready)
		},
	}

	go func() {
		err := config.Start(serveCtx, r.server)
		r.serveDone <- err
		close(r.serveDone)
	}()

	select {
	case <-ready:
		return nil
	case err := <-r.serveDone:
		cancel()
		if err != nil {
			return fmt.Errorf("start Echo server on %s: %w", r.address, err)
		}
		return errors.New("echo server stopped before becoming ready")
	case <-ctx.Done():
		cancel()
		return ctx.Err()
	}
}

// Shutdown cancels Echo's serving context, allowing StartConfig to perform
// graceful HTTP shutdown, and waits for the serving goroutine to finish.
func (r *echoServerResource) Shutdown(ctx context.Context) error {
	r.shutdownOnce.Do(func() {
		if r.cancel != nil {
			r.cancel()
		}
		go func() {
			var err error
			if r.serveDone != nil {
				if serveErr, ok := <-r.serveDone; ok {
					err = serveErr
				}
			}
			r.shutdownErr = err
			close(r.shutdownDone)
		}()
	})

	select {
	case <-r.shutdownDone:
		return r.shutdownErr
	case <-ctx.Done():
		return ctx.Err()
	}
}
