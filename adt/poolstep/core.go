package poolstep

import (
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/poolexp"
	"orglang/go-engine/adt/poolvar"
	"orglang/go-engine/adt/procdec"
	"orglang/go-engine/adt/procexec"
)

type API interface {
}

type StepSpec struct {
	ImplRef implsem.SemRef
	PoolExp poolexp.ExpSpec
}

// step (aka Sem)
type StepRec interface {
	comm()
}

// publication (aka msg)
type PubRec struct {
	// совпадает со значением в SubRec
	CommRef commsem.SemRef
	ImplRef implsem.SemRef
	ChnlID  identity.ADT
	ValExp  poolexp.ExpRec
}

func (r PubRec) comm() {}

// subscription (aka srv)
type SubRec struct {
	// совпадает со значением в PubRec
	CommRef commsem.SemRef
	ImplRef implsem.SemRef
	ChnlID  identity.ADT
	ContExp poolexp.ExpRec
}

func (r SubRec) comm() {}

type service struct {
	poolSteps Repo
	implSems  implsem.Repo
	procExecs procexec.Repo
	poolVars  poolvar.Repo
	procDecs  procdec.Repo
	operator  db.Operator
	log       *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	poolSteps Repo,
	implSems implsem.Repo,
	procExecs procexec.Repo,
	poolVars poolvar.Repo,
	procDecs procdec.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{poolSteps, implSems, procExecs, poolVars, procDecs, operator, log.With(name)}
}

func ErrRecTypeUnexpected(got StepRec) error {
	return fmt.Errorf("step rec unexpected: %T%+v", got, got)
}

func ErrRecTypeMismatch(got, want StepRec) error {
	return fmt.Errorf("step rec mismatch: want %T, got %T", want, got)
}

func ErrStepKindUnexpected(got stepKind) error {
	return fmt.Errorf("step kind unexpected: %v", got)
}
