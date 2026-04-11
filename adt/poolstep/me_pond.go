package poolstep

import (
	"log/slog"
	"reflect"

	"github.com/alitto/pond/v2"
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

func (b *pondBroker) SendRec(rec StepRec) error {
	panic("unimplemented")
}

func (b *pondBroker) SendSpec(spec StepSpec) error {
	panic("unimplemented")
}
