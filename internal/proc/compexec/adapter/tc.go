package adapter

import (
	"orglang/go-engine/internal/adt/compsem/adapter"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
	compexec "orglang/go-engine/internal/proc/compexec/core"
)

func DataFromMod(source compexec.ExecMod) (ExecMod, error) {
	compRefs, err := adapter.DataFromRefs(source.CompRefs)
	if err != nil {
		return ExecMod{}, err
	}
	return ExecMod{
		CompRefs:   compRefs,
		LinearVars: compvaradapter.DataFromLinearRecs(source.LinearVars),
	}, nil
}
