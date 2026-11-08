package adapter

import (
	"orglang/go-engine/internal/adt/connsem"
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	termexp "orglang/go-engine/internal/pool/termexp/adapter"

	commturn "orglang/go-engine/internal/pool/commturn/core"
)

func DataFromStepRec(r commturn.TurnRec) commturn.TurnRecDS {
	switch rec := r.(type) {
	case commturn.PubRec:
		connRef := connsem.DataFromRef(rec.ConnRef)
		compRef := compsem.DataFromRef(rec.CompRef)
		return commturn.TurnRecDS{
			ConnID: connRef.ConnID,
			ConnRN: connRef.ConnRN,
			CompID: compRef.CompID,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commturn.PubKind,
			Exp:    termexp.DataFromExpRec(rec.ValExp),
		}
	case commturn.SubRec:
		connRef := connsem.DataFromRef(rec.ConnRef)
		compRef := compsem.DataFromRef(rec.CompRef)
		return commturn.TurnRecDS{
			ConnID: connRef.ConnID,
			ConnRN: connRef.ConnRN,
			CompID: compRef.CompID,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commturn.SubKind,
			Exp:    termexp.DataFromExpRec(rec.ContExp),
		}
	default:
		panic(commturn.ErrRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto commturn.TurnRecDS) (commturn.TurnRec, error) {
	switch dto.K {
	case commturn.PubKind:
		return DataToPubRec(dto)
	case commturn.SubKind:
		return DataToSubRec(dto)
	default:
		panic(commturn.ErrStepKindUnexpected(dto.K))
	}
}
