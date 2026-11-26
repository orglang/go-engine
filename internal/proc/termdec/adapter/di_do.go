package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/proc/ent"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

var Package = do.Package(
	do.Lazy[termdec.Repo](func(i do.Injector) (termdec.Repo, error) {
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoEnt(client, logger), nil
	}),
)
