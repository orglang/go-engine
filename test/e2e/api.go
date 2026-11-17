package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/orglang/go-engine/internal/lib/e2e/api"
	"github.com/orglang/go-sdk/adt/compsem"
	"github.com/orglang/go-sdk/adt/termsem"
	"github.com/orglang/go-sdk/adt/termvar"
	"github.com/orglang/go-sdk/adt/typesem"
	poolcompexec "github.com/orglang/go-sdk/pool/compexec"
	poolcompstep "github.com/orglang/go-sdk/pool/compstep"
	pooltermdec "github.com/orglang/go-sdk/pool/termdec"
	pooltypedef "github.com/orglang/go-sdk/pool/typedef"
	proccompstep "github.com/orglang/go-sdk/proc/compstep"
	proctermdec "github.com/orglang/go-sdk/proc/termdec"
	proctypedef "github.com/orglang/go-sdk/proc/typedef"
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

func contractJSON(src, dst any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	renameContractKeys(value)
	data, err = json.Marshal(value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func renameContractKeys(value any) {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			newKey := key
			switch key {
			case "chnl_ph", "comm_ph":
				newKey = "conn_ph"
			case "cont_exp", "cont_es":
				newKey = "cont_expr"
			case "pool_exp":
				newKey = "pool_expr"
			case "proc_exp":
				newKey = "proc_expr"
			case "type_exp":
				newKey = "type_expr"
			case "exp":
				newKey = "expr"
			case "val_es":
				newKey = "val_expr"
			case "lab_qn":
				newKey = "pattern_qn"
			}
			if newKey != key {
				delete(v, key)
				v[newKey] = child
			}
			renameContractKeys(child)
		}
	case []any:
		for _, child := range v {
			renameContractKeys(child)
		}
	}
}

func responseError(status int, body []byte) error {
	if status >= http.StatusBadRequest {
		return fmt.Errorf("generated API returned %d: %s", status, body)
	}
	return nil
}

type poolDeclarationClient struct{ api *generatedAPI }

func (c *poolDeclarationClient) Create(spec pooltermdec.DecSpec) (typesem.SemRef, error) {
	var body api.PoolDecSpec
	if err := contractJSON(spec, &body); err != nil { return typesem.SemRef{}, err }
	r, err := c.api.client.CreatePoolDeclarationWithResponse(context.Background(), body)
	if err != nil { return typesem.SemRef{}, err }
	if err := responseError(r.StatusCode(), r.GetBody()); err != nil { return typesem.SemRef{}, err }
	ref := r.GetJSON201()
	if ref == nil { return typesem.SemRef{}, fmt.Errorf("create pool declaration: empty response") }
	return typesem.SemRef{TypeID: ref.TermId, TypeRN: ref.TermRn}, nil
}

type poolExecutionClient struct{ api *generatedAPI }

func (c *poolExecutionClient) Create(spec poolcompexec.ExecSpec) (compsem.SemRef, error) {
	var body api.AnyExecSpec
	if err := contractJSON(spec, &body); err != nil { return compsem.SemRef{}, err }
	r, err := c.api.client.CreatePoolExecutionWithResponse(context.Background(), body)
	if err != nil { return compsem.SemRef{}, err }
	if err := responseError(r.StatusCode(), r.GetBody()); err != nil { return compsem.SemRef{}, err }
	ref := r.GetJSON201()
	if ref == nil { return compsem.SemRef{}, fmt.Errorf("create pool execution: empty response") }
	return compsem.SemRef{CompID: ref.CompId, CompRN: ref.CompRn}, nil
}

func (c *poolExecutionClient) Take(spec poolcompstep.StepSpec) error {
	var body api.PoolStepSpec
	if err := contractJSON(spec, &body); err != nil { return err }
	r, err := c.api.client.StepPoolExecutionWithResponse(context.Background(), body)
	if err != nil { return err }
	return responseError(r.StatusCode(), r.GetBody())
}

func (c *poolExecutionClient) Spawn(spec poolcompstep.StepSpec) (compsem.SemRef, error) {
	var body api.PoolStepSpec
	if err := contractJSON(spec, &body); err != nil { return compsem.SemRef{}, err }
	r, err := c.api.client.SpawnPoolExecutionWithResponse(context.Background(), body)
	if err != nil { return compsem.SemRef{}, err }
	if err := responseError(r.StatusCode(), r.GetBody()); err != nil { return compsem.SemRef{}, err }
	ref := r.GetJSON201()
	if ref == nil { return compsem.SemRef{}, fmt.Errorf("spawn pool execution: empty response") }
	return compsem.SemRef{CompID: ref.CompId, CompRN: ref.CompRn}, nil
}

type poolTypeDefinitionClient struct{ api *generatedAPI }

func (c *poolTypeDefinitionClient) Create(spec pooltypedef.DefSpec) (typesem.SemRef, error) {
	var body api.PoolTypeDefSpec
	if err := contractJSON(spec, &body); err != nil { return typesem.SemRef{}, err }
	r, err := c.api.client.CreatePoolTypeDefinitionWithResponse(context.Background(), body)
	if err != nil { return typesem.SemRef{}, err }
	if err := responseError(r.StatusCode(), r.GetBody()); err != nil { return typesem.SemRef{}, err }
	ref := r.GetJSON201()
	if ref == nil { return typesem.SemRef{}, fmt.Errorf("create pool type definition: empty response") }
	return typesem.SemRef{TypeID: ref.TypeId, TypeRN: ref.TypeRn}, nil
}

type processDeclarationClient struct{ api *generatedAPI }

func (c *processDeclarationClient) Create(spec proctermdec.DecSpec) (proctermdec.DecSnap, error) {
	var body api.ProcDecSpec
	if err := contractJSON(spec, &body); err != nil { return proctermdec.DecSnap{}, err }
	r, err := c.api.client.CreateProcessDeclarationWithResponse(context.Background(), body)
	if err != nil { return proctermdec.DecSnap{}, err }
	if err := responseError(r.StatusCode(), r.GetBody()); err != nil { return proctermdec.DecSnap{}, err }
	ref := r.GetJSON201()
	if ref == nil { return proctermdec.DecSnap{}, fmt.Errorf("create process declaration: empty response") }
	snap := proctermdec.DecSnap{TermRef: termsem.SemRef{TermID: ref.TermId, TermRN: ref.TermRn}}
	snap.LiabVar = termvar.VarRec{ChnlPH: spec.LiabVar.ChnlPH}
	snap.AssetVars = make([]termvar.VarRec, len(spec.AssetVars))
	for i, variable := range spec.AssetVars {
		snap.AssetVars[i] = termvar.VarRec{ChnlPH: variable.ChnlPH}
	}
	return snap, nil
}

type processExecutionClient struct{ api *generatedAPI }

func (c *processExecutionClient) Take(spec proccompstep.StepSpec) error {
	var body api.ProcStepSpec
	if err := contractJSON(spec, &body); err != nil { return err }
	r, err := c.api.client.StepProcessExecutionWithResponse(context.Background(), spec.CompRef.CompID, body)
	if err != nil { return err }
	return responseError(r.StatusCode(), r.GetBody())
}

type processTypeDefinitionClient struct{ api *generatedAPI }

func (c *processTypeDefinitionClient) Create(spec proctypedef.DefSpec) (proctypedef.DefSnap, error) {
	var body api.ProcTypeDefSpec
	if err := contractJSON(spec, &body); err != nil { return proctypedef.DefSnap{}, err }
	r, err := c.api.client.CreateProcessTypeDefinitionWithResponse(context.Background(), body)
	if err != nil { return proctypedef.DefSnap{}, err }
	if err := responseError(r.StatusCode(), r.GetBody()); err != nil { return proctypedef.DefSnap{}, err }
	ref := r.GetJSON201()
	if ref == nil { return proctypedef.DefSnap{}, fmt.Errorf("create process type definition: empty response") }
	return proctypedef.DefSnap{
		TypeRef: typesem.SemRef{TypeID: ref.TypeId, TypeRN: ref.TypeRn},
		DefSpec: spec,
	}, nil
}

// Keep these compile-time assertions close to the adapters.
var (
	_ interface{ Create(pooltermdec.DecSpec) (typesem.SemRef, error) } = (*poolDeclarationClient)(nil)
	_ interface{ Create(poolcompexec.ExecSpec) (compsem.SemRef, error); Take(poolcompstep.StepSpec) error; Spawn(poolcompstep.StepSpec) (compsem.SemRef, error) } = (*poolExecutionClient)(nil)
	_ interface{ Create(pooltypedef.DefSpec) (typesem.SemRef, error) } = (*poolTypeDefinitionClient)(nil)
	_ interface{ Create(proctermdec.DecSpec) (proctermdec.DecSnap, error) } = (*processDeclarationClient)(nil)
	_ interface{ Take(proccompstep.StepSpec) error } = (*processExecutionClient)(nil)
	_ interface{ Create(proctypedef.DefSpec) (proctypedef.DefSnap, error) } = (*processTypeDefinitionClient)(nil)
)
