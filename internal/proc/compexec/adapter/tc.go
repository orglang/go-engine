package adapter

import (
	"orglang/go-engine/internal/adt/comproot/adapter"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
	compexec "orglang/go-engine/internal/proc/compexec/core"
)

func DataFromMod(source compexec.EvalMod) (EvalMod, error) {
	compRefs, err := adapter.DataFromRefs(source.CompRefs)
	if err != nil {
		return EvalMod{}, err
	}
	return EvalMod{
		CompRefs:   compRefs,
		LinearVars: compvaradapter.DataFromLinearRecs(source.LinearVars),
	}, nil
}
