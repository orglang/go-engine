package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/alitto/pond/v2"

	compstep "orglang/go-engine/internal/pool/compstep/core"
	execeval "orglang/go-engine/internal/pool/execeval/core"
)

type pondBroker struct {
	api     execeval.API
	pool    pond.Pool
	taskCtx context.Context
	log     *slog.Logger
}

func newPondBroker(taskCtx context.Context, pool pond.Pool, log *slog.Logger) *pondBroker {
	name := slog.String("name", reflect.TypeFor[pondBroker]().Name())
	return &pondBroker{nil, pool, taskCtx, log.With(name)}
}

func cfgPondBroker(broker execeval.Broker, api execeval.API) {
	broker.Subscribe(api)
}

// for compilation purposes
func newPondExch() execeval.Broker {
	return new(pondBroker)
}

func (b *pondBroker) Subscribe(api execeval.API) {
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
