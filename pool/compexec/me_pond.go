package compexec

import (
	"log/slog"
	"orglang/go-engine/pool/compstep"
	"reflect"

	"github.com/alitto/pond/v2"
)

type pondBroker struct {
	api  API
	pool pond.Pool
	log  *slog.Logger
}

func newPondBroker(pool pond.Pool, log *slog.Logger) *pondBroker {
	name := slog.String("name", reflect.TypeFor[pondBroker]().Name())
	return &pondBroker{nil, pool, log.With(name)}
}

func cfgPondBroker(broker Broker, api API) error {
	broker.Subscribe(api)
	return nil
}

// for compilation purposes
func newPondExch() Broker {
	return new(pondBroker)
}

func (b *pondBroker) Subscribe(api API) {
	b.api = api
}

func (b *pondBroker) SendSpec(spec compstep.StepSpec) error {
	b.pool.Go(func() {
		apiErr := b.api.Take(spec)
		if apiErr != nil {
			b.log.Error("consumption failed", slog.Any("ref", spec.CompRef), slog.Any("reason", apiErr))
			return
		}
		b.log.Debug("consumption succeed", slog.Any("ref", spec.CompRef))
	})
	return nil
}
