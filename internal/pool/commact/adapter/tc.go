package adapter

import (
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	termexp "orglang/go-engine/internal/pool/termexp/adapter"

	commact "orglang/go-engine/internal/pool/commact/core"
)

func DataFromStepRec(r commact.ActRec) commact.ActRecDS {
	switch rec := r.(type) {
	case commact.PubRec:
		compRef := compsem.DataFromRef(rec.CompRef)
		return commact.ActRecDS{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commact.PubKind,
			Exp:    termexp.DataFromExpRec(rec.ValExp),
		}
	case commact.SubRec:
		compRef := compsem.DataFromRef(rec.CompRef)
		return commact.ActRecDS{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commact.SubKind,
			Exp:    termexp.DataFromExpRec(rec.ContExp),
		}
	default:
		panic(commact.ErrActRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto commact.ActRecDS) (commact.ActRec, error) {
	switch dto.K {
	case commact.PubKind:
		return DataToPubRec(dto)
	case commact.SubKind:
		return DataToSubRec(dto)
	default:
		panic(commact.ErrActKindUnexpected(dto.K))
	}
}
