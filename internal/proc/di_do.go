package proc

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descsem"
	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"
	"orglang/go-engine/internal/adt/typesem"
	typesemadapter "orglang/go-engine/internal/adt/typesem/adapter"
	compconn "orglang/go-engine/internal/proc/compconn/adapter"
	compexec "orglang/go-engine/internal/proc/compexec/adapter"
	termdec "orglang/go-engine/internal/proc/termdec/adapter"
	termdef "orglang/go-engine/internal/proc/termdef/adapter"
	typedef "orglang/go-engine/internal/proc/typedef/adapter"
	typeexp "orglang/go-engine/internal/proc/typeexp/adapter"
)

// Package is the process-domain composition boundary. It registers the
// persistence adapters and services currently used by production process
// endpoints. Incomplete and unused areas such as termexp remain unwired.
var Package = do.Package(
	compconn.Package,
	compexec.Package,
	termdec.Package,
	termdef.Package,
	typedef.Package,
	typeexp.Package,
	do.Lazy[typesem.Repo](func(i do.Injector) (typesem.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return typesemadapter.NewDaoPgx("proc_type_defs ", "proc_desc_binds ")(logger), nil
	}),
	do.Lazy[descsem.Repo](func(i do.Injector) (descsem.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return descsemadapter.NewDaoPgx("proc_desc_binds ")(logger), nil
	}),
)
