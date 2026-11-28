package execvar

import (
	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/identity"
)

func SampleVarRec() LinearRec {
	return LinearRec{
		RootRef: comproot.New(),
		ConnID:  identity.New(),
		ChnlID:  identity.New(),
		ChnlPH:  "chnl-1",
		ExpCK:   1,
		ChnlBS:  LiabSide,
	}
}
