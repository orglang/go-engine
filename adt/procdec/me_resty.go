package procdec

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
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

func (cl *sdkResty) Incept(decQN qualsym.ADT) (DecRef, error) {
	return DecRef{}, nil
}

func (cl *sdkResty) Create(spec DecSpec) (DecSnap, error) {
	req := MsgFromDecSpec(spec)
	var res DecSnapME
	resp, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/declarations")
	if err != nil {
		return DecSnap{}, err
	}
	if resp.IsError() {
		return DecSnap{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToDecSnap(res)
}

func (c *sdkResty) RetrieveSnap(id identity.ADT) (DecSnap, error) {
	var res DecSnapME
	resp, err := c.resty.R().
		SetResult(&res).
		SetPathParam("id", id.String()).
		Get("/declarations/{id}")
	if err != nil {
		return DecSnap{}, err
	}
	if resp.IsError() {
		return DecSnap{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToDecSnap(res)
}

func (c *sdkResty) RetreiveRefs() ([]DecRef, error) {
	refs := []DecRef{}
	return refs, nil
}
