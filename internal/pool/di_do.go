package pool

import (
	"github.com/samber/do/v2"

	compconn "orglang/go-engine/internal/pool/compconn/adapter"
	compvar "orglang/go-engine/internal/pool/compvar/adapter"
	typeexp "orglang/go-engine/internal/pool/typeexp/adapter"
)

// Package is the pool-domain composition boundary. The initial migrated slice
// contains the persistence adapters that have no cross-domain service wiring.
var Package = do.Package(
	compconn.Package,
	compvar.Package,
	typeexp.Package,
)
