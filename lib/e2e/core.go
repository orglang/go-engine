package e2e

import (
	"github.com/go-resty/resty/v2"

	"github.com/orglang/go-sdk/adt/descsem"
	"github.com/orglang/go-sdk/adt/implsem"
	"github.com/orglang/go-sdk/adt/pooldec"
	"github.com/orglang/go-sdk/adt/poolexec"
	"github.com/orglang/go-sdk/adt/poolstep"
	"github.com/orglang/go-sdk/adt/procdec"
	"github.com/orglang/go-sdk/adt/procexec"
	"github.com/orglang/go-sdk/adt/procstep"
	"github.com/orglang/go-sdk/adt/typedef"
	"github.com/orglang/go-sdk/adt/xactdef"
)

type PoolDecAPI interface {
	Create(pooldec.DecSpec) (descsem.SemRef, error)
}

func newPoolDecAPI(client *resty.Client) PoolDecAPI {
	return &pooldec.RestySDK{Client: client}
}

type PoolExecAPI interface {
	Retrieve(implsem.SemRef) (poolexec.ExecSnap, error)
	Create(poolexec.ExecSpec) (implsem.SemRef, error)
	Take(poolstep.StepSpec) error
	Spawn(poolstep.StepSpec) (implsem.SemRef, error)
	Poll(poolexec.PollSpec) (implsem.SemRef, error)
}

func newPoolExecAPI(client *resty.Client) PoolExecAPI {
	return &poolexec.RestySDK{Client: client}
}

type XactDefAPI interface {
	Create(xactdef.DefSpec) (descsem.SemRef, error)
}

func newXactDefAPI(client *resty.Client) XactDefAPI {
	return &xactdef.RestySDK{Client: client}
}

type ProcDecAPI interface {
	Create(procdec.DecSpec) (procdec.DecSnap, error)
}

func newProcDecAPI(client *resty.Client) ProcDecAPI {
	return &procdec.RestySDK{Client: client}
}

type ProcExecAPI interface {
	Take(procstep.StepSpec) error
}

func newProcExecAPI(client *resty.Client) ProcExecAPI {
	return &procexec.RestySDK{Client: client}
}

type TypeDefAPI interface {
	Create(typedef.DefSpec) (typedef.DefSnap, error)
}

func newTypeDefAPI(client *resty.Client) TypeDefAPI {
	return &typedef.RestySDK{Client: client}
}
