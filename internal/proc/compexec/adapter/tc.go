package adapter

import (
	"orglang/go-engine/internal/adt/compsem/adapter"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
	compexec "orglang/go-engine/internal/proc/compexec/core"
)

func DataFromMod(source compexec.ExecMod) (ExecModDS, error) {
	compRefs, err := adapter.DataFromRefs(source.CompRefs)
	if err != nil {
		return ExecModDS{}, err
	}
	return ExecModDS{
		CompRefs:   compRefs,
		LinearVars: compvaradapter.DataFromLinearRecs(source.LinearVars),
	}, nil
}
