package pool

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descsem"
	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"
	"orglang/go-engine/internal/lib/wp"
	compconn "orglang/go-engine/internal/pool/compconn/adapter"
	compexec "orglang/go-engine/internal/pool/compexec/adapter"
	compvar "orglang/go-engine/internal/pool/compvar/adapter"
	termdef "orglang/go-engine/internal/pool/termdef/adapter"
	typedef "orglang/go-engine/internal/pool/typedef/adapter"
	typeexp "orglang/go-engine/internal/pool/typeexp/adapter"
	proccompexec "orglang/go-engine/internal/proc/compexec/adapter"
)

// Package is the pool-domain composition boundary. Persistence adapters and
// services are registered here; inbound HTTP route registration is explicit
// and remains a separate composition concern.
var Package = do.Package(
	wp.Package,
	proccompexec.Package,
	compconn.Package,
	compvar.Package,
	compexec.Package,
	typeexp.Package,
	typedef.Package,
	termdef.Package,
	do.Lazy[descsem.Repo](func(i do.Injector) (descsem.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return descsemadapter.NewDaoPgx("pool_desc_binds ")(logger), nil
	}),
)
