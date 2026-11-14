// Package program provides a concise API for submitting programs.
package program

import (
	"context"
	"strings"

	ogen "orglang/go-engine/prototypes/api-client/ogen"
)

// Client exposes program operations without generator-prefixed names.
type Client struct {
	api *ogen.Client
}

// New wraps a configured ogen client.
func New(api *ogen.Client) *Client { return &Client{api: api} }

// Create submits source text for parsing and registration.
func (c *Client) Create(ctx context.Context, source string) error {
	return c.api.EngineApiCreateProgram(ctx, ogen.EngineApiCreateProgramReq{
		Data: strings.NewReader(source),
	})
}
