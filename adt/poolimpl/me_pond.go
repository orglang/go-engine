package poolimpl

import (
	"log/slog"
	"reflect"

	"github.com/alitto/pond/v2"

	"orglang/go-engine/adt/poolstep"
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

func cfgPondBroker(broker Exch, api API) error {
	broker.Subscribe(api)
	return nil
}

// for compilation purposes
func newPondExch() Exch {
	return new(pondBroker)
}

func (b *pondBroker) Subscribe(api API) {
	b.api = api
}

func (b *pondBroker) SendSpec(spec poolstep.StepSpec) error {
	b.pool.Go(func() {
		apiErr := b.api.Take(spec)
		if apiErr != nil {
			b.log.Error("consumption failed", slog.Any("ref", spec.ImplRef), slog.Any("reason", apiErr))
			return
		}
		b.log.Debug("consumption succeed", slog.Any("ref", spec.ImplRef))
	})
	return nil
}
