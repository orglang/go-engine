package core

import (
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/commsem"
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	compvar "orglang/go-engine/internal/pool/compvar/core"
	termexp "orglang/go-engine/internal/pool/termexp/core"
	proccompexec "orglang/go-engine/internal/proc/compexec/core"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

type API any

// step (aka Sem)
type TurnRec interface {
	turn()
}

// publication (aka msg)
type PubRec struct {
	// совпадает со значением в SubRec
	CommRef commsem.SemRef
	CompRef compsem.SemRef
	ChnlID  identity.ADT
	ValExp  termexp.ExpRec
}

func (r PubRec) turn() {}

// subscription (aka srv)
type SubRec struct {
	// совпадает со значением в PubRec
	CommRef commsem.SemRef
	CompRef compsem.SemRef
	ChnlID  identity.ADT
	ContExp termexp.ExpRec
}

func (r SubRec) turn() {}

type service struct {
	poolSteps  Repo
	procExecs  proccompexec.Repo
	poolVars   compvar.Repo
	procDecs   termdec.Repo
	transactor tx.Transactor
	log        *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func NewService(
	poolSteps Repo,
	procExecs proccompexec.Repo,
	poolVars compvar.Repo,
	procDecs termdec.Repo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{poolSteps, procExecs, poolVars, procDecs, transactor, log.With(name)}
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
