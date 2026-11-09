package adapter

import (
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	termexp "orglang/go-engine/internal/pool/termexp/adapter"

	connact "orglang/go-engine/internal/pool/connact/core"
)

func DataFromStepRec(r connact.ActRec) connact.ActRecDS {
	switch rec := r.(type) {
	case connact.PubRec:
		compRef := compsem.DataFromRef(rec.CompRef)
		return connact.ActRecDS{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      connact.PubKind,
			Exp:    termexp.DataFromExpRec(rec.ValExp),
		}
	case connact.SubRec:
		compRef := compsem.DataFromRef(rec.CompRef)
		return connact.ActRecDS{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      connact.SubKind,
			Exp:    termexp.DataFromExpRec(rec.ContExp),
		}
	default:
		panic(connact.ErrActRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto connact.ActRecDS) (connact.ActRec, error) {
	switch dto.K {
	case connact.PubKind:
		return DataToPubRec(dto)
	case connact.SubKind:
		return DataToSubRec(dto)
	default:
		panic(connact.ErrActKindUnexpected(dto.K))
	}
}
