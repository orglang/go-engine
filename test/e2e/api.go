package e2e

import (
	"context"
	"fmt"
	"net/http"

	api "orglang/go-engine/internal/lib/e2e/api"
)

type generatedAPI struct {
	client *api.ClientWithResponses
}

func newGeneratedAPI(server string) (*generatedAPI, error) {
	client, err := api.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}
	return &generatedAPI{client: client}, nil
}

func responseError(status int, body []byte) error {
	if status >= http.StatusBadRequest {
		return fmt.Errorf("generated API returned %d: %s", status, body)
	}
	return nil
}

func (c *generatedAPI) createPoolDeclaration(spec api.PoolDecSpec) (api.AnyTermRef, error) {
	r, err := c.client.CreatePoolDeclarationWithResponse(context.Background(), spec)
	if err != nil {
		return api.AnyTermRef{}, err
	}
	if err := responseError(r.StatusCode(), r.GetBody()); err != nil {
		return api.AnyTermRef{}, err
	}
	if ref := r.GetJSON201(); ref != nil {
		return *ref, nil
	}
	return api.AnyTermRef{}, fmt.Errorf("create pool declaration: empty response")
}

func (c *generatedAPI) createPoolExecution(spec api.AnyExecSpec) (api.AnyCompRef, error) {
	r, err := c.client.CreatePoolExecutionWithResponse(context.Background(), spec)
	if err != nil {
		return api.AnyCompRef{}, err
	}
	if err := responseError(r.StatusCode(), r.GetBody()); err != nil {
		return api.AnyCompRef{}, err
	}
	if ref := r.GetJSON201(); ref != nil {
		return *ref, nil
	}
	return api.AnyCompRef{}, fmt.Errorf("create pool execution: empty response")
}

func (c *generatedAPI) stepPoolExecution(spec api.PoolStepSpec) error {
	r, err := c.client.StepPoolExecutionWithResponse(context.Background(), spec)
	if err != nil {
		return err
	}
	return responseError(r.StatusCode(), r.GetBody())
}

func (c *generatedAPI) createPoolTypeDefinition(spec api.PoolTypeDefSpec) (api.AnyTypeRef, error) {
	r, err := c.client.CreatePoolTypeDefinitionWithResponse(context.Background(), spec)
	if err != nil {
		return api.AnyTypeRef{}, err
	}
	if err := responseError(r.StatusCode(), r.GetBody()); err != nil {
		return api.AnyTypeRef{}, err
	}
	if ref := r.GetJSON201(); ref != nil {
		return *ref, nil
	}
	return api.AnyTypeRef{}, fmt.Errorf("create pool type definition: empty response")
}
