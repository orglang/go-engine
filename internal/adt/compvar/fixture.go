package compvar

import (
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/connsem"
	"orglang/go-engine/internal/adt/identity"
)

func SampleVarRec() LinearRec {
	return LinearRec{
		CompRef: compsem.New(),
		ConnRef: connsem.New(),
		ChnlID:  identity.New(),
		ChnlPH:  "chnl-1",
		ExpVK:   1,
		ChnlBS:  LiabSide,
	}
}
