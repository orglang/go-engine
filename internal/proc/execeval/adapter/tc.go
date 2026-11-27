package adapter

import (
	"orglang/go-engine/internal/adt/comproot/adapter"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
	execeval "orglang/go-engine/internal/proc/execeval/core"
)

func DataFromMod(source execeval.EvalMod) (EvalMod, error) {
	compRefs, err := adapter.DataFromRefs(source.CompRefs)
	if err != nil {
		return EvalMod{}, err
	}
	return EvalMod{
		CompRefs:   compRefs,
		LinearVars: compvaradapter.DataFromLinearRecs(source.LinearVars),
	}, nil
}
