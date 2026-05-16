package commturn

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/pool/termexp"
)

func DataFromStepRec(r TurnRec) TurnRecDS {
	switch rec := r.(type) {
	case PubRec:
		commRef := commsem.DataFromRef(rec.CommRef)
		compRef := compsem.DataFromRef(rec.CompRef)
		return TurnRecDS{
			CommID: commRef.CommID,
			CommRN: commRef.CommRN,
			CompID: compRef.CompID,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      PubKind,
			Exp:    termexp.DataFromExpRec(rec.ValExp),
		}
	case SubRec:
		commRef := commsem.DataFromRef(rec.CommRef)
		compRef := compsem.DataFromRef(rec.CompRef)
		return TurnRecDS{
			CommID: commRef.CommID,
			CommRN: commRef.CommRN,
			CompID: compRef.CompID,
			K:      SubKind,
			Exp:    termexp.DataFromExpRec(rec.ContExp),
		}
	default:
		panic(ErrRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto TurnRecDS) (TurnRec, error) {
	switch dto.K {
	case PubKind:
		return DataToPubRec(dto)
	case SubKind:
		return DataToSubRec(dto)
	default:
		panic(ErrStepKindUnexpected(dto.K))
	}
}
