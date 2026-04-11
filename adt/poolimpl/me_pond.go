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

func newPondBroker(api API, pool pond.Pool, log *slog.Logger) *pondBroker {
	name := slog.String("name", reflect.TypeFor[pondBroker]().Name())
	return &pondBroker{api, pool, log.With(name)}
}

// for compilation purposes
func newExch() Exch {
	return new(pondBroker)
}

func (b *pondBroker) SendSpec(spec poolstep.StepSpec) error {
	b.pool.SubmitErr(func() error {
		apiErr := b.api.Take(spec)
		if apiErr != nil {
			return apiErr
		}
		return nil
	})
	return nil
}
