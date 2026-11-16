package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/app"
	"orglang/go-engine/internal/app/web"
	"orglang/go-engine/internal/lib/db"
	"orglang/go-engine/internal/lib/kv"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/wp"
	"orglang/go-engine/internal/lib/ws"
	"orglang/go-engine/internal/pool"
	"orglang/go-engine/internal/proc"
)

const shutdownTimeout = 30 * time.Second

func main() {
	if err := run(); err != nil {
		slog.Error("engine stopped with an error", slog.Any("reason", err))
		os.Exit(1)
	}
}

func run() error {
	application, err := newApp()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Start(ctx); err != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return errors.Join(fmt.Errorf("start application: %w", err), application.Stop(shutdownCtx))
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := application.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("stop application: %w", err)
	}
	return nil
}

func newApp() (*app.App, error) {
	application := app.NewApp(
		db.Package,
		kv.Package,
		lf.Package,
		wp.Package,
		ws.Package,
		pool.Package,
		proc.Package,
		web.Package,
	)

	server, err := do.Invoke[*echo.Echo](application.Container())
	if err != nil {
		shutdownUnresolvedApp(application)
		return nil, fmt.Errorf("resolve HTTP server: %w", err)
	}
	controller, err := do.Invoke[web.Controller](application.Container())
	if err != nil {
		shutdownUnresolvedApp(application)
		return nil, fmt.Errorf("resolve web controller: %w", err)
	}
	web.RegisterRoutes(server, controller)

	return application, nil
}

func shutdownUnresolvedApp(application *app.App) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := application.Stop(ctx); err != nil {
		slog.Error("failed to clean up application after composition error", slog.Any("reason", err))
	}
}
