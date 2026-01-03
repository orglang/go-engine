package poolxec

import (
	"github.com/go-resty/resty/v2"

	"orglang/orglang/adt/identity"

	"orglang/orglang/adt/procxec"
)

// Client-side secondary adapter
type sdkResty struct {
	resty *resty.Client
}

func newSdkResty() *sdkResty {
	r := resty.New().SetBaseURL("http://localhost:8080/api/v1")
	return &sdkResty{r}
}

func NewAPI() API {
	return newSdkResty()
}

func (cl *sdkResty) Create(spec ExecSpec) (ExecRef, error) {
	req := MsgFromExecSpec(spec)
	var res ExecRefME
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/pools")
	if err != nil {
		return ExecRef{}, err
	}
	return MsgToExecRef(res)
}

func (cl *sdkResty) Poll(spec PollSpec) (procxec.ExecRef, error) {
	return procxec.ExecRef{}, nil
}

func (cl *sdkResty) Retrieve(poolID identity.ADT) (ExecSnap, error) {
	var res ExecSnapME
	_, err := cl.resty.R().
		SetResult(&res).
		SetPathParam("id", poolID.String()).
		Get("/pools/{id}")
	if err != nil {
		return ExecSnap{}, err
	}
	return MsgToExecSnap(res)
}

func (cl *sdkResty) RetreiveRefs() ([]ExecRef, error) {
	refs := []ExecRef{}
	return refs, nil
}

func (cl *sdkResty) Spawn(spec procxec.ExecSpec) (procxec.ExecRef, error) {
	req := procxec.MsgFromExecSpec(spec)
	var res procxec.ExecRefME
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		SetPathParam("poolID", spec.PoolID.String()).
		Post("/pools/{poolID}/procs")
	if err != nil {
		return procxec.ExecRef{}, err
	}
	return procxec.MsgToExecRef(res)
}

func (cl *sdkResty) Take(spec StepSpec) error {
	req := MsgFromStepSpec(spec)
	var res procxec.ExecRefME
	_, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		SetPathParam("poolID", spec.PoolID.String()).
		SetPathParam("procID", spec.ProcID.String()).
		Post("/pools/{poolID}/procs/{procID}/steps")
	if err != nil {
		return err
	}
	return nil
}
