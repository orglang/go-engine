package adapter

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/identity"
	termexp "orglang/go-engine/pool/termexp/adapter"

	commturn "orglang/go-engine/pool/commturn/core"
)

func DataFromStepRec(r commturn.TurnRec) commturn.TurnRecDS {
	switch rec := r.(type) {
	case commturn.PubRec:
		commRef := commsem.DataFromRef(rec.CommRef)
		compRef := compsem.DataFromRef(rec.CompRef)
		return commturn.TurnRecDS{
			CommID: commRef.CommID,
			CommRN: commRef.CommRN,
			CompID: compRef.CompID,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commturn.PubKind,
			Exp:    termexp.DataFromExpRec(rec.ValExp),
		}
	case commturn.SubRec:
		commRef := commsem.DataFromRef(rec.CommRef)
		compRef := compsem.DataFromRef(rec.CompRef)
		return commturn.TurnRecDS{
			CommID: commRef.CommID,
			CommRN: commRef.CommRN,
			CompID: compRef.CompID,
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
