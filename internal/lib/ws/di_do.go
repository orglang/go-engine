package ws

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
			server: server,
			address: fmt.Sprintf(":%d", config.Protocol.HTTP.Port),
			shutdownDone: make(chan struct{}),
		}
		if err := lifecycle.Append(app.Hook{
			OnStart: resource.start,
			OnStop: resource.Shutdown,
		}); err != nil {
			return nil, err
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
	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:   true,
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(_ echo.Context, values middleware.RequestLoggerValues) error {
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
	server *echo.Echo
	address string

	mu sync.Mutex
	listener net.Listener
	serveDone chan error
	shutdownOnce sync.Once
	shutdownDone chan struct{}
	shutdownErr error
}

func (r *echoServerResource) start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	listener, err := net.Listen("tcp", r.address)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", r.address, err)
	}

	r.mu.Lock()
	r.listener = listener
	r.serveDone = make(chan error, 1)
	serveDone := r.serveDone
	r.mu.Unlock()

	go func() {
		err := r.server.StartListener(listener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serveDone <- err
		close(serveDone)
	}()

	return nil
}

// Shutdown gracefully stops Echo and waits for its serving goroutine.
func (r *echoServerResource) Shutdown(ctx context.Context) error {
	r.shutdownOnce.Do(func() {
		go func() {
			err := r.server.Shutdown(ctx)
			if err != nil {
				err = errors.Join(err, r.server.Close())
			}

			r.mu.Lock()
			serveDone := r.serveDone
			r.mu.Unlock()
			if serveDone != nil {
				if serveErr, ok := <-serveDone; ok {
					err = errors.Join(err, serveErr)
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
