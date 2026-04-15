package compexec

import (
	"log/slog"
	"orglang/go-engine/pool/compstep"
	"reflect"

	"github.com/gammazero/workerpool"
)

type workerPoolBroker struct {
	api  API
	pool *workerpool.WorkerPool
	log  *slog.Logger
}

func newWorkerPoolBroker(pool *workerpool.WorkerPool, log *slog.Logger) *workerPoolBroker {
	name := slog.String("name", reflect.TypeFor[workerPoolBroker]().Name())
	return &workerPoolBroker{nil, pool, log.With(name)}
}

func cfgWorkerPoolBroker(broker Exch, api API) error {
	broker.Subscribe(api)
	return nil
}

// for compilation purposes
func newWorkerPoolExch() Exch {
	return new(workerPoolBroker)
}

func (b *workerPoolBroker) Subscribe(api API) {
	b.api = api
}

func (b *workerPoolBroker) SendSpec(spec compstep.StepSpec) error {
	b.pool.Submit(func() {
		apiErr := b.api.Take(spec)
		if apiErr != nil {
			b.log.Error("consumption failed", slog.Any("ref", spec.CompRef), slog.Any("reason", apiErr))
			return
		}
		b.log.Debug("consumption succeed", slog.Any("ref", spec.CompRef))
	})
	return nil
}
