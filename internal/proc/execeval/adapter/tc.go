package adapter

import (
	"orglang/go-engine/internal/adt/execroot/adapter"
	execvaradapter "orglang/go-engine/internal/adt/execvar/adapter"
	execeval "orglang/go-engine/internal/proc/execeval/core"
)

func DataFromMod(source execeval.EvalMod) (EvalMod, error) {
	rootRefs, err := adapter.DataFromRefs(source.RootRefs)
	if err != nil {
		return EvalMod{}, err
	}
	return EvalMod{
		RootRefs:   rootRefs,
		LinearVars: execvaradapter.DataFromLinearRecs(source.LinearVars),
	}, nil
}
