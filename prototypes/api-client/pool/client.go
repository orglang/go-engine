// Package pool provides a domain-oriented, concise API over the generated ogen client.
package pool

import (
	"context"

	ogen "orglang/go-engine/prototypes/api-client/ogen"
)

// Client exposes pool operations without generator- or service-prefixed names.
type Client struct {
	api *ogen.Client
}

// New wraps a configured ogen client.
func New(api *ogen.Client) *Client { return &Client{api: api} }

type (
	TypeDefSpec = ogen.PoolTypeDefSpec
	TypeDefSnap = ogen.PoolTypeDefSnap
	TypeExpr = ogen.PoolTypeExprSpec
	TypeOne = ogen.PoolTypeOneSpec
	TypeLink = ogen.PoolTypeLinkSpec
	TypeTensor = ogen.PoolTypeTensorSpec
	TypeLolli = ogen.PoolTypeLolliSpec
	TypePlus = ogen.PoolTypePlusSpec
	TypeWith = ogen.PoolTypeWithSpec
	TypeUp = ogen.PoolTypeUpSpec
	TypeDown = ogen.PoolTypeDownSpec
	TypeRef = ogen.AnyTypeRef
	TermExpr = ogen.PoolTermExprSpec
	TermRef = ogen.AnyTermRef
)

// One constructs the "one" pool type expression variant.
func One() TypeExpr { return ogen.NewPoolTypeOneSpecPoolTypeExprSpec(TypeOne{}) }

// Link constructs the "link" pool type expression variant.
func Link(v TypeLink) TypeExpr { return ogen.NewPoolTypeLinkSpecPoolTypeExprSpec(v) }

// Tensor constructs the "tensor" pool type expression variant.
func Tensor(v TypeTensor) TypeExpr { return ogen.NewPoolTypeTensorSpecPoolTypeExprSpec(v) }

// Lolli constructs the "lolli" pool type expression variant.
func Lolli(v TypeLolli) TypeExpr { return ogen.NewPoolTypeLolliSpecPoolTypeExprSpec(v) }

// Plus constructs the "plus" pool type expression variant.
func Plus(v TypePlus) TypeExpr { return ogen.NewPoolTypePlusSpecPoolTypeExprSpec(v) }

// With constructs the "with" pool type expression variant.
func With(v TypeWith) TypeExpr { return ogen.NewPoolTypeWithSpecPoolTypeExprSpec(v) }

// Up constructs the "up" pool type expression variant.
func Up(v TypeUp) TypeExpr { return ogen.NewPoolTypeUpSpecPoolTypeExprSpec(v) }

// Down constructs the "down" pool type expression variant.
func Down(v TypeDown) TypeExpr { return ogen.NewPoolTypeDownSpecPoolTypeExprSpec(v) }

// CreateType creates a pool type definition.
func (c *Client) CreateType(ctx context.Context, spec *TypeDefSpec) (*TypeRef, error) {
	return c.api.EngineApiCreatePoolTypeDefinition(ctx, spec)
}

// GetType returns the latest pool type definition snapshot.
func (c *Client) GetType(ctx context.Context, id string) (*TypeDefSnap, error) {
	return c.api.EngineApiGetPoolTypeDefinition(ctx, ogen.EngineApiGetPoolTypeDefinitionParams{ID: id})
}

// GetTypeAt returns a pool type definition snapshot at the requested revision.
func (c *Client) GetTypeAt(ctx context.Context, id string, revision int64) (*TypeDefSnap, error) {
	return c.api.EngineApiGetPoolTypeDefinition(ctx, ogen.EngineApiGetPoolTypeDefinitionParams{
		ID: id, Rn: ogen.NewOptInt64(revision),
	})
}
