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
	do.Lazy[*taskContextResource](func(do.Injector) (*taskContextResource, error) {
		taskCtx, cancel := context.WithCancel(context.Background())
		return &taskContextResource{ctx: taskCtx, cancel: cancel}, nil
	}),
	do.Lazy[context.Context](func(i do.Injector) (context.Context, error) {
		resource, err := do.Invoke[*taskContextResource](i)
		if err != nil {
			return nil, err
		}
		return resource.ctx, nil
	}),
	do.Lazy[*pondPoolResource](func(i do.Injector) (*pondPoolResource, error) {
		resource, err := do.Invoke[*taskContextResource](i)
		if err != nil {
			return nil, err
		}
		return &pondPoolResource{
			pool:   newPondPool(),
			cancel: resource.cancel,
			done:   make(chan struct{}),
		}, nil
	}),
	do.Lazy[pond.Pool](func(i do.Injector) (pond.Pool, error) {
		resource, err := do.Invoke[*pondPoolResource](i)
		if err != nil {
			return nil, err
		}
		return resource.pool, nil
	}),
	do.Lazy[*workerPoolResource](func(i do.Injector) (*workerPoolResource, error) {
		resource, err := do.Invoke[*taskContextResource](i)
		if err != nil {
			return nil, err
		}
		return &workerPoolResource{
			pool:   newWorkerPool(),
			cancel: resource.cancel,
			done:   make(chan struct{}),
		}, nil
	}),
	do.Lazy[*workerpool.WorkerPool](func(i do.Injector) (*workerpool.WorkerPool, error) {
		resource, err := do.Invoke[*workerPoolResource](i)
		if err != nil {
			return nil, err
		}
		return resource.pool, nil
	}),
)

type taskContextResource struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type pondPoolResource struct {
	pool   pond.Pool
	cancel context.CancelFunc
	once   sync.Once
	done   chan struct{}
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

type workerPoolResource struct {
	pool   *workerpool.WorkerPool
	cancel context.CancelFunc
	once   sync.Once
	done   chan struct{}
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
