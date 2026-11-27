package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/gammazero/workerpool"

	execeval "orglang/go-engine/internal/pool/execeval/core"
	compstep "orglang/go-engine/internal/pool/compstep/core"
)

type workerPoolBroker struct {
	api     execeval.API
	pool    *workerpool.WorkerPool
	taskCtx context.Context
	log     *slog.Logger
}

func newWorkerPoolBroker(taskCtx context.Context, pool *workerpool.WorkerPool, log *slog.Logger) *workerPoolBroker {
	name := slog.String("name", reflect.TypeFor[workerPoolBroker]().Name())
	return &workerPoolBroker{nil, pool, taskCtx, log.With(name)}
}

func cfgWorkerPoolBroker(broker execeval.Broker, api execeval.API) {
	broker.Subscribe(api)
}

// for compilation purposes
func newWorkerPoolExch() execeval.Broker {
	return new(workerPoolBroker)
}

func (b *workerPoolBroker) Subscribe(api execeval.API) {
	b.api = api
}

// SendSpec starts work owned by the worker pool, not by the originating HTTP request.
func (b *workerPoolBroker) SendSpec(spec compstep.StepSpec) error {
	b.pool.Submit(func() {
		apiErr := b.api.Take(b.taskCtx, spec)
		if apiErr != nil {
			b.log.Error("consumption failed", slog.Any("ref", spec.CompRef), slog.Any("reason", apiErr))
			return
		}
		b.log.Debug("consumption succeed", slog.Any("ref", spec.CompRef))
	})
	return nil
}
