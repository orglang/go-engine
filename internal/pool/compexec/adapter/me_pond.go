package adapter

import (
	"log/slog"
	"reflect"

	"github.com/alitto/pond/v2"

	compexec "orglang/go-engine/internal/pool/compexec/core"
	compstep "orglang/go-engine/internal/pool/compstep/core"
	"orglang/go-engine/internal/lib/wp"
)

type pondBroker struct {
	api  compexec.API
	pool pond.Pool
	taskCtx context.Context
	log  *slog.Logger
}

func newPondBroker(pool pond.Pool, taskCtx wp.PondTaskContext, log *slog.Logger) *pondBroker {
	name := slog.String("name", reflect.TypeFor[pondBroker]().Name())
	return &pondBroker{nil, pool, taskCtx, log.With(name)}
}

func cfgPondBroker(broker compexec.Broker, api compexec.API) {
	broker.Subscribe(api)
}

// for compilation purposes
func newPondExch() compexec.Broker {
	return new(pondBroker)
}

func (b *pondBroker) Subscribe(api compexec.API) {
	b.api = api
}

// SendSpec starts work owned by the worker pool, not by the originating HTTP request.
func (b *pondBroker) SendSpec(spec compstep.StepSpec) error {
	return b.pool.Go(func() {
		apiErr := b.api.Take(b.taskCtx, spec)
		if apiErr != nil {
			b.log.Error("consumption failed", slog.Any("ref", spec.CompRef), slog.Any("reason", apiErr))
			return
		}
		b.log.Debug("consumption succeed", slog.Any("ref", spec.CompRef))
	})
}
