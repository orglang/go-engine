package compvar

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/identity"
)

func SampleVarRec() LinearRec {
	return LinearRec{
		CompRef: compsem.New(),
		CommRef: commsem.New(),
		ChnlID:  identity.New(),
		ChnlPH:  "chnl-1",
		ExpVK:   1,
		ChnlBS:  LiabSide,
	}
}
