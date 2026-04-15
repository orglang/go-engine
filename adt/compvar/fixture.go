package compvar

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/semcomm"
	"orglang/go-engine/adt/semterm"
)

func SampleVarRec() LinearRec {
	return LinearRec{
		TermRef: semterm.NewRef(),
		ExchRef: semcomm.NewRef(),
		ChnlID:  identity.New(),
		ChnlPH:  "chnl-1",
		ExpVK:   1,
		ChnlBS:  LiabSide,
	}
}
