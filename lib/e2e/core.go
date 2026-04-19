package e2e

import (
	"github.com/go-resty/resty/v2"

	"github.com/orglang/go-sdk/adt/compsem"
	"github.com/orglang/go-sdk/adt/typesem"
	"github.com/orglang/go-sdk/pool/compexec"
	"github.com/orglang/go-sdk/pool/compstep"
	termdec1 "github.com/orglang/go-sdk/pool/termdec"
	typedef1 "github.com/orglang/go-sdk/pool/typedef"
	compexec1 "github.com/orglang/go-sdk/proc/compexec"
	compstep1 "github.com/orglang/go-sdk/proc/compstep"
	"github.com/orglang/go-sdk/proc/termdec"
	"github.com/orglang/go-sdk/proc/typedef"
)

type PoolDecAPI interface {
	Create(termdec1.DecSpec) (typesem.SemRef, error)
}

func newPoolDecAPI(client *resty.Client) PoolDecAPI {
	return &termdec1.RestySDK{Client: client}
}

type PoolExecAPI interface {
	Retrieve(compsem.SemRef) (compexec.ExecSnap, error)
	Create(compexec.ExecSpec) (compsem.SemRef, error)
	Take(compstep.StepSpec) error
	Spawn(compstep.StepSpec) (compsem.SemRef, error)
}

func newPoolExecAPI(client *resty.Client) PoolExecAPI {
	return &compexec.RestySDK{Client: client}
}

type XactDefAPI interface {
	Create(typedef1.DefSpec) (typesem.SemRef, error)
}

func newXactDefAPI(client *resty.Client) XactDefAPI {
	return &typedef1.RestySDK{Client: client}
}

type ProcDecAPI interface {
	Create(termdec.DecSpec) (termdec.DecSnap, error)
}

func newProcDecAPI(client *resty.Client) ProcDecAPI {
	return &termdec.RestySDK{Client: client}
}

type ProcExecAPI interface {
	Take(compstep1.StepSpec) error
}

func newProcExecAPI(client *resty.Client) ProcExecAPI {
	return &compexec1.RestySDK{Client: client}
}

type TypeDefAPI interface {
	Create(typedef.DefSpec) (typedef.DefSnap, error)
}

func newTypeDefAPI(client *resty.Client) TypeDefAPI {
	return &typedef.RestySDK{Client: client}
}
