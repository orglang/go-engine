package typedef

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

func (cl *sdkResty) Incept(typeQN qualsym.ADT) (DefRef, error) {
	return DefRef{}, nil
}

func (cl *sdkResty) Create(spec DefSpec) (DefSnap, error) {
	req := MsgFromDefSpec(spec)
	var res DefSnapME
	resp, err := cl.resty.R().
		SetResult(&res).
		SetBody(&req).
		Post("/types")
	if err != nil {
		return DefSnap{}, err
	}
	if resp.IsError() {
		return DefSnap{}, fmt.Errorf("received: %v", string(resp.Body()))
	}
	return MsgToDefSnap(res)
}

func (c *sdkResty) Modify(snap DefSnap) (DefSnap, error) {
	return DefSnap{}, nil
}

func (c *sdkResty) Retrieve(defID identity.ADT) (DefSnap, error) {
	return DefSnap{}, nil
}

func (c *sdkResty) retrieveSnap(rec DefRec) (DefSnap, error) {
	return DefSnap{}, nil
}

func (c *sdkResty) RetreiveRefs() ([]DefRef, error) {
	return []DefRef{}, nil
}
