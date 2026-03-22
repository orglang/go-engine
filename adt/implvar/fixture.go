package implvar

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
)

func SampleVarRec() LinearRec {
	return LinearRec{
		ImplRef: implsem.NewRef(),
		CommRef: commsem.NewRef(),
		ChnlID:  identity.New(),
		ChnlPH:  "chnl-1",
		ExpVK:   1,
		ChnlBS:  LiabSide,
	}
}
