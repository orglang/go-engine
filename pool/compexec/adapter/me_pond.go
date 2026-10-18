package adapter

import (
	"log/slog"
	"reflect"

	"github.com/alitto/pond/v2"

	compexec "orglang/go-engine/pool/compexec/core"
	compstep "orglang/go-engine/pool/compstep/core"
)

type pondBroker struct {
	api  compexec.API
	pool pond.Pool
	log  *slog.Logger
}

func newPondBroker(pool pond.Pool, log *slog.Logger) *pondBroker {
	name := slog.String("name", reflect.TypeFor[pondBroker]().Name())
	return &pondBroker{nil, pool, log.With(name)}
}

func cfgPondBroker(broker compexec.Broker, api compexec.API) error {
	broker.Subscribe(api)
	return nil
}

// for compilation purposes
func newPondExch() compexec.Broker {
	return new(pondBroker)
}

func (b *pondBroker) Subscribe(api compexec.API) {
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
