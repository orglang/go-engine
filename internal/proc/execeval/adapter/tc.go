package adapter

import (
	"orglang/go-engine/internal/adt/comproot/adapter"
	execvaradapter "orglang/go-engine/internal/adt/execvar/adapter"
	execeval "orglang/go-engine/internal/proc/execeval/core"
)

func DataFromMod(source execeval.EvalMod) (EvalMod, error) {
	compRefs, err := adapter.DataFromRefs(source.CompRefs)
	if err != nil {
		return EvalMod{}, err
	}
	return EvalMod{
		CompRefs:   compRefs,
		LinearVars: execvaradapter.DataFromLinearRecs(source.LinearVars),
	}, nil
}
