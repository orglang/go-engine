package proc

import (
	"github.com/samber/do/v2"

	compconn "orglang/go-engine/internal/proc/compconn/adapter"
	compexec "orglang/go-engine/internal/proc/compexec/adapter"
	termdec "orglang/go-engine/internal/proc/termdec/adapter"
	termdef "orglang/go-engine/internal/proc/termdef/adapter"
	typedef "orglang/go-engine/internal/proc/typedef/adapter"
	typeexp "orglang/go-engine/internal/proc/typeexp/adapter"
)

// Package is the process-domain composition boundary. It registers the
// persistence adapters currently used by production process services.
// Incomplete and unused areas such as termexp remain intentionally unwired.
var Package = do.Package(
	compconn.Package,
	compexec.Package,
	termdec.Package,
	termdef.Package,
	typedef.Package,
	typeexp.Package,
)
