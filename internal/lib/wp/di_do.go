package wp

import (
	"context"
	"sync"

	"github.com/alitto/pond/v2"
	"github.com/gammazero/workerpool"
	"github.com/samber/do/v2"
)

// Package registers worker pools and gives the application container explicit
// ownership of their shutdown.
var Package = do.Package(
	do.Lazy[*pondPoolResource](func(do.Injector) (*pondPoolResource, error) {
		taskCtx, cancel := context.WithCancel(context.Background())
		return &pondPoolResource{
			pool:    newPondPool(),
			taskCtx: taskCtx,
			cancel:  cancel,
			done: make(chan struct{}),
		}, nil
	}),
	do.Lazy[PondTaskContext](func(i do.Injector) (PondTaskContext, error) {
		resource, err := do.Invoke[*pondPoolResource](i)
		if err != nil {
			return PondTaskContext{}, err
		}
		return PondTaskContext{Context: resource.taskCtx}, nil
	}),
	do.Lazy[pond.Pool](func(i do.Injector) (pond.Pool, error) {
		resource, err := do.Invoke[*pondPoolResource](i)
		if err != nil {
			return nil, err
		}
		return resource.pool, nil
	}),
	do.Lazy[*workerPoolResource](func(do.Injector) (*workerPoolResource, error) {
		taskCtx, cancel := context.WithCancel(context.Background())
		return &workerPoolResource{
			pool:    newWorkerPool(),
			taskCtx: taskCtx,
			cancel:  cancel,
			done: make(chan struct{}),
		}, nil
	}),
	do.Lazy[WorkerPoolTaskContext](func(i do.Injector) (WorkerPoolTaskContext, error) {
		resource, err := do.Invoke[*workerPoolResource](i)
		if err != nil {
			return WorkerPoolTaskContext{}, err
		}
		return WorkerPoolTaskContext{Context: resource.taskCtx}, nil
	}),
	do.Lazy[*workerpool.WorkerPool](func(i do.Injector) (*workerpool.WorkerPool, error) {
		resource, err := do.Invoke[*workerPoolResource](i)
		if err != nil {
			return nil, err
		}
		return resource.pool, nil
	}),
)

// PondTaskContext is canceled when the application begins shutting down.
type PondTaskContext struct {
	context.Context
}

type pondPoolResource struct {
	pool    pond.Pool
	taskCtx context.Context
	cancel  context.CancelFunc
	once    sync.Once
	done    chan struct{}
}

// Shutdown stops accepting work and waits for all accepted pond tasks.
func (r *pondPoolResource) Shutdown(ctx context.Context) error {
	r.once.Do(func() {
		r.cancel()
		go func() {
			r.pool.StopAndWait()
			close(r.done)
		}()
	})

	select {
	case <-r.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WorkerPoolTaskContext is canceled when the application begins shutting down.
type WorkerPoolTaskContext struct {
	context.Context
}

type workerPoolResource struct {
	pool    *workerpool.WorkerPool
	taskCtx context.Context
	cancel  context.CancelFunc
	once    sync.Once
	done    chan struct{}
}

// Shutdown cancels the worker task context and drains accepted tasks before returning.
func (r *workerPoolResource) Shutdown(ctx context.Context) error {
	r.once.Do(func() {
		r.cancel()
		go func() {
			r.pool.StopWait()
			close(r.done)
		}()
	})

	select {
	case <-r.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
