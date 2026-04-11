package poolstep

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/poolexp"
)

func DataFromStepRec(r StepRec) StepRecDS {
	switch rec := r.(type) {
	case PubRec:
		commRef := commsem.DataFromRef(rec.CommRef)
		implRef := implsem.DataFromRef(rec.ImplRef)
		return StepRecDS{
			CommID: commRef.CommID,
			CommRN: commRef.CommRN,
			ImplID: implRef.ImplID,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      PubKind,
			Exp:    poolexp.DataFromExpRec(rec.ValExp),
		}
	case SubRec:
		commRef := commsem.DataFromRef(rec.CommRef)
		implRef := implsem.DataFromRef(rec.ImplRef)
		return StepRecDS{
			CommID: commRef.CommID,
			CommRN: commRef.CommRN,
			ImplID: implRef.ImplID,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      SubKind,
			Exp:    poolexp.DataFromExpRec(rec.ContExp),
		}
	default:
		panic(ErrRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto StepRecDS) (StepRec, error) {
	switch dto.K {
	case PubKind:
		commRef, err := DataToCommRef(dto)
		if err != nil {
			return nil, err
		}
		implRef, err := DataToImplRef(dto)
		if err != nil {
			return nil, err
		}
		chnlID, err := identity.ConvertFromString(dto.ChnlID)
		if err != nil {
			return nil, err
		}
		valExp, err := poolexp.DataToExpRec(dto.Exp)
		if err != nil {
			return nil, err
		}
		return PubRec{CommRef: commRef, ImplRef: implRef, ChnlID: chnlID, ValExp: valExp}, nil
	case SubKind:
		commRef, err := DataToCommRef(dto)
		if err != nil {
			return nil, err
		}
		implRef, err := DataToImplRef(dto)
		if err != nil {
			return nil, err
		}
		chnlID, err := identity.ConvertFromString(dto.ChnlID)
		if err != nil {
			return nil, err
		}
		contExp, err := poolexp.DataToExpRec(dto.Exp)
		if err != nil {
			return nil, err
		}
		return SubRec{CommRef: commRef, ImplRef: implRef, ChnlID: chnlID, ContExp: contExp}, nil
	default:
		panic(ErrStepKindUnexpected(dto.K))
	}
}
