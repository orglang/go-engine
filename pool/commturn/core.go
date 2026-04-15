package commturn

import (
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/semcomm"
	"orglang/go-engine/adt/semcomp"
	"orglang/go-engine/pool/compvar"
	"orglang/go-engine/pool/termexp"
	"orglang/go-engine/proc/compexec"
	"orglang/go-engine/proc/termdec"
)

type API interface {
}

// step (aka Sem)
type TurnRec interface {
	turn()
}

// publication (aka msg)
type PubRec struct {
	// совпадает со значением в SubRec
	CommRef semcomm.CommRef
	CompRef semcomp.CompRef
	ChnlID  identity.ADT
	ValExp  termexp.ExpRec
}

func (r PubRec) turn() {}

// subscription (aka srv)
type SubRec struct {
	// совпадает со значением в PubRec
	CommRef semcomm.CommRef
	CompRef semcomp.CompRef
	ChnlID  identity.ADT
	ContExp termexp.ExpRec
}

func (r SubRec) turn() {}

type service struct {
	poolSteps Repo
	procExecs compexec.Repo
	poolVars  compvar.Repo
	procDecs  termdec.Repo
	operator  db.Operator
	log       *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	poolSteps Repo,
	procExecs compexec.Repo,
	poolVars compvar.Repo,
	procDecs termdec.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{poolSteps, procExecs, poolVars, procDecs, operator, log.With(name)}
}

func ErrRecTypeUnexpected(got TurnRec) error {
	return fmt.Errorf("step rec unexpected: %T%+v", got, got)
}

func ErrRecTypeMismatch(got, want TurnRec) error {
	return fmt.Errorf("step rec mismatch: want %T, got %T", want, got)
}

func ErrStepKindUnexpected(got turnKind) error {
	return fmt.Errorf("step kind unexpected: %v", got)
}
