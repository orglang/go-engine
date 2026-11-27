package adapter

import (
	"fmt"
	comprootadapter "orglang/go-engine/internal/adt/comproot/adapter"
	"orglang/go-engine/internal/adt/identity"
	compexec "orglang/go-engine/internal/pool/compexec/adapter"

	connact "orglang/go-engine/internal/pool/connact/core"
)

func DataFromStepRec(r connact.ActRec) ActRec {
	switch rec := r.(type) {
	case connact.PubRec:
		compRef := comprootadapter.DataFromRef(rec.CompRef)
		return ActRec{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      PubKind,
			Exp:    compexec.DataFromExpRec(rec.ValExp),
		}
	case connact.SubRec:
		compRef := comprootadapter.DataFromRef(rec.CompRef)
		return ActRec{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      SubKind,
			Exp:    compexec.DataFromExpRec(rec.ContExp),
		}
	default:
		panic(connact.ErrActRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto ActRec) (connact.ActRec, error) {
	switch dto.K {
	case PubKind:
		return DataToPubRec(dto)
	case SubKind:
		return DataToSubRec(dto)
	default:
		panic(fmt.Errorf("act kind unexpected: %v", dto.K))
	}
}

