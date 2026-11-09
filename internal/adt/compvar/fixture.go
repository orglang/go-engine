package compvar

import (
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
)

func SampleVarRec() LinearRec {
	return LinearRec{
		CompRef: compsem.New(),
		ConnID:  identity.New(),
		ChnlID:  identity.New(),
		ChnlPH:  "chnl-1",
		ExpCK:   1,
		ChnlBS:  LiabSide,
	}
}
