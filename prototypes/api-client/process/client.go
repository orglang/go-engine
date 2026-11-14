// Package process provides a domain-oriented, concise API over the generated ogen client.
package process

import (
	"context"

	ogen "orglang/go-engine/prototypes/api-client/ogen"
)

// Client exposes process operations without generator- or service-prefixed names.
type Client struct {
	api *ogen.Client
}

// New wraps a configured ogen client.
func New(api *ogen.Client) *Client { return &Client{api: api} }

type (
	TypeDefSpec = ogen.ProcTypeDefSpec
	TypeDefSnap = ogen.ProcTypeDefSnap
	TypeExpr = ogen.ProcTypeExprSpec
	TypeOne = ogen.ProcTypeOneSpec
	TypeLink = ogen.ProcTypeLinkSpec
	TypeTensor = ogen.ProcTypeTensorSpec
	TypeLolli = ogen.ProcTypeLolliSpec
	TypePlus = ogen.ProcTypePlusSpec
	TypeWith = ogen.ProcTypeWithSpec
	TypeUp = ogen.ProcTypeUpSpec
	TypeDown = ogen.ProcTypeDownSpec
	TypeRef = ogen.AnyTypeRef
	DeclarationSpec = ogen.ProcDecSpec
	DeclarationSnap = ogen.ProcDecSnap
	TermExpr = ogen.ProcTermExprSpec
	TermRef = ogen.AnyTermRef
	ExecutionSnap = ogen.ProcCompSnap
	StepSpec = ogen.ProcStepSpec
)

// One constructs the "one" process type expression variant.
func One() TypeExpr { return ogen.NewProcTypeOneSpecProcTypeExprSpec(TypeOne{}) }

// Link constructs the "link" process type expression variant.
func Link(v TypeLink) TypeExpr { return ogen.NewProcTypeLinkSpecProcTypeExprSpec(v) }

// Tensor constructs the "tensor" process type expression variant.
func Tensor(v TypeTensor) TypeExpr { return ogen.NewProcTypeTensorSpecProcTypeExprSpec(v) }

// Lolli constructs the "lolli" process type expression variant.
func Lolli(v TypeLolli) TypeExpr { return ogen.NewProcTypeLolliSpecProcTypeExprSpec(v) }

// Plus constructs the "plus" process type expression variant.
func Plus(v TypePlus) TypeExpr { return ogen.NewProcTypePlusSpecProcTypeExprSpec(v) }

// With constructs the "with" process type expression variant.
func With(v TypeWith) TypeExpr { return ogen.NewProcTypeWithSpecProcTypeExprSpec(v) }

// Up constructs the "up" process type expression variant.
func Up(v TypeUp) TypeExpr { return ogen.NewProcTypeUpSpecProcTypeExprSpec(v) }

// Down constructs the "down" process type expression variant.
func Down(v TypeDown) TypeExpr { return ogen.NewProcTypeDownSpecProcTypeExprSpec(v) }

// CreateType creates a process type definition.
func (c *Client) CreateType(ctx context.Context, spec *TypeDefSpec) (*TypeRef, error) {
	return c.api.EngineApiCreateProcessTypeDefinition(ctx, spec)
}

// GetType returns the latest process type definition snapshot.
func (c *Client) GetType(ctx context.Context, id string) (*TypeDefSnap, error) {
	return c.api.EngineApiGetProcessTypeDefinition(ctx, ogen.EngineApiGetProcessTypeDefinitionParams{ID: id})
}

// GetTypeAt returns a process type definition snapshot at the requested revision.
func (c *Client) GetTypeAt(ctx context.Context, id string, revision int64) (*TypeDefSnap, error) {
	return c.api.EngineApiGetProcessTypeDefinition(ctx, ogen.EngineApiGetProcessTypeDefinitionParams{
		ID: id, Rn: ogen.NewOptInt64(revision),
	})
}

// CreateDeclaration creates a process declaration.
func (c *Client) CreateDeclaration(ctx context.Context, spec *DeclarationSpec) (*TermRef, error) {
	return c.api.EngineApiCreateProcessDeclaration(ctx, spec)
}

// GetDeclaration returns a process declaration at the latest revision.
func (c *Client) GetDeclaration(ctx context.Context, id string) (*DeclarationSnap, error) {
	return c.api.EngineApiGetProcessDeclaration(ctx, ogen.EngineApiGetProcessDeclarationParams{ID: id})
}

// GetDefinition returns the latest process term definition.
func (c *Client) GetDefinition(ctx context.Context, id string) (TermExpr, error) {
	return c.api.EngineApiGetProcessDefinition(ctx, ogen.EngineApiGetProcessDefinitionParams{ID: id})
}

// ExtendDefinition adds expression parts to a process term definition.
func (c *Client) ExtendDefinition(ctx context.Context, id string, extension TermExpr) (*TermRef, error) {
	return c.api.EngineApiExtendProcessDefinition(ctx, extension, ogen.EngineApiExtendProcessDefinitionParams{ID: id})
}

// GetExecution returns a process execution at the latest available revision.
func (c *Client) GetExecution(ctx context.Context, id string) (*ExecutionSnap, error) {
	return c.api.EngineApiGetProcessExecution(ctx, ogen.EngineApiGetProcessExecutionParams{ID: id})
}

// StepExecution applies one step to a process execution.
func (c *Client) StepExecution(ctx context.Context, id string, spec *StepSpec) error {
	return c.api.EngineApiStepProcessExecution(ctx, spec, ogen.EngineApiStepProcessExecutionParams{ID: id})
}

// ListExecutions lists execution identities associated with a process.
func (c *Client) ListExecutions(ctx context.Context, id string) ([]ogen.AnyCompRef, error) {
	return c.api.EngineApiListProcessExecutions(ctx, ogen.EngineApiListProcessExecutionsParams{ID: id})
}
