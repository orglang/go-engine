package execvar

import (
	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/identity"
)

func SampleVarRec() LinearRec {
	return LinearRec{
		ExecRef: comproot.New(),
		ConnID:  identity.New(),
		ChnlID:  identity.New(),
		ChnlPH:  "chnl-1",
		ExpCK:   1,
		ChnlBS:  LiabSide,
	}
}
