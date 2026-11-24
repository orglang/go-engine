// Package core is a compatibility facade for the unified process termdef capability.
package core

import (
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
	termdef "orglang/go-engine/internal/proc/termdef/core"
	termexp "orglang/go-engine/internal/proc/termexp/core"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type API = termdef.API
type DecSpec = termdef.DecSpec
type DecRec = termdef.DecRec
type DecSnap = termdef.DecSnap
type Repo = termdef.Repo

func NewService(
	termDecRepo Repo,
	termExpRepo termexp.Repo,
	descRootRepo descroot.Repo,
	descSemRepo descsem.Repo,
	typeDefRepo typedef.Repo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	return termdef.NewService(termDecRepo, termExpRepo, descRootRepo, descSemRepo, typeDefRepo, transactor, log)
}

func CollectEnv(recs iter.Seq[DecRec]) []uniqsym.ADT {
	return termdef.CollectEnv(recs)
}

func ErrRootMissingInEnv(rid identity.ADT) error {
	return fmt.Errorf("root missing in env: %v", rid)
}
