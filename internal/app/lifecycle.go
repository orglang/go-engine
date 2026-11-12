package app

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/samber/do/v2"
)

// Hook describes one application resource's startup and shutdown operations.
// OnStop is called only when OnStart completed successfully.
type Hook struct {
	OnStart func(context.Context) error
	OnStop  func(context.Context) error
}

// Lifecycle registers and runs application resource hooks.
type Lifecycle interface {
	Append(Hook) error
	Start(context.Context) error
	Stop(context.Context) error
}

type lifecycleState uint8

const (
	lifecycleNew lifecycleState = iota
	lifecycleStarting
	lifecycleStarted
	lifecycleStopping
	lifecycleStopped
)

type lifecycle struct {
	mu      sync.Mutex
	state   lifecycleState
	hooks   []Hook
	started []Hook
}

// NewLifecycle creates an independent lifecycle coordinator.
func NewLifecycle() Lifecycle {
	return &lifecycle{}
}

func (l *lifecycle) Append(hook Hook) error {
	if hook.OnStart == nil && hook.OnStop == nil {
		return errors.New("lifecycle hook must define OnStart or OnStop")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.state != lifecycleNew {
		return errors.New("cannot append lifecycle hook after startup has begun")
	}
	l.hooks = append(l.hooks, hook)
	return nil
}

// Start runs hooks in registration order. If a hook fails, already-started
// hooks are stopped in reverse order and the startup and rollback errors are joined.
func (l *lifecycle) Start(ctx context.Context) error {
	l.mu.Lock()
	if l.state != lifecycleNew {
		state := l.state
		l.mu.Unlock()
		return fmt.Errorf("cannot start application lifecycle in state %d", state)
	}
	l.state = lifecycleStarting
	hooks := append([]Hook(nil), l.hooks...)
	l.mu.Unlock()

	started := make([]Hook, 0, len(hooks))
	for _, hook := range hooks {
		if hook.OnStart != nil {
			if err := hook.OnStart(ctx); err != nil {
				errs := []error{fmt.Errorf("start lifecycle hook: %w", err)}
				for i := len(started) - 1; i >= 0; i-- {
					if started[i].OnStop != nil {
						if stopErr := started[i].OnStop(ctx); stopErr != nil {
							errs = append(errs, fmt.Errorf("rollback lifecycle hook: %w", stopErr))
						}
					}
				}
				l.mu.Lock()
				l.state = lifecycleStopped
				l.mu.Unlock()
				return errors.Join(errs...)
			}
		}
		started = append(started, hook)
	}

	l.mu.Lock()
	l.started = started
	l.state = lifecycleStarted
	l.mu.Unlock()
	return nil
}

// Stop runs successfully started hooks in reverse registration order.
// Repeated calls are harmless.
func (l *lifecycle) Stop(ctx context.Context) error {
	l.mu.Lock()
	switch l.state {
	case lifecycleStopped:
		l.mu.Unlock()
		return nil
	case lifecycleNew:
		l.state = lifecycleStopped
		l.mu.Unlock()
		return nil
	case lifecycleStarted:
		l.state = lifecycleStopping
	default:
		state := l.state
		l.mu.Unlock()
		return fmt.Errorf("cannot stop application lifecycle in state %d", state)
	}
	started := append([]Hook(nil), l.started...)
	l.mu.Unlock()

	var errs []error
	for i := len(started) - 1; i >= 0; i-- {
		if started[i].OnStop != nil {
			if err := started[i].OnStop(ctx); err != nil {
				errs = append(errs, fmt.Errorf("stop lifecycle hook: %w", err))
			}
		}
	}

	l.mu.Lock()
	l.state = lifecycleStopped
	l.mu.Unlock()
	return errors.Join(errs...)
}

// App owns a DI container and the lifecycle of resources registered with it.
// Services that implement samber/do's Shutdowner interfaces are also shut down
// by the container when App.Stop is called.
type App struct {
	injector  do.Injector
	lifecycle Lifecycle
}

// NewApp creates an application with an isolated DI container and lifecycle.
// Providers can resolve Lifecycle through do.Invoke[Lifecycle]. Lazy providers
// are not instantiated automatically; composition code must resolve required
// services before calling Start.
func NewApp(packages ...func(do.Injector)) *App {
	injector := NewContainer(packages...)
	lifecycle := NewLifecycle()
	do.ProvideValue[Lifecycle](injector, lifecycle)
	return &App{injector: injector, lifecycle: lifecycle}
}

// Container returns the application's dependency-injection container.
func (a *App) Container() do.Injector {
	return a.injector
}

// Start starts application resources in registration order.
func (a *App) Start(ctx context.Context) error {
	return a.lifecycle.Start(ctx)
}

// Stop stops lifecycle hooks in reverse order, then shuts down instantiated
// services managed by samber/do. Both phases are attempted even if one fails.
func (a *App) Stop(ctx context.Context) error {
	hookErr := a.lifecycle.Stop(ctx)
	report := a.injector.ShutdownWithContext(ctx)
	if report == nil {
		return hookErr
	}
	return errors.Join(hookErr, report)
}
