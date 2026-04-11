package poolstep

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/poolexp"
)

func DataFromStepRec(r StepRec) StepRecDS {
	switch rec := r.(type) {
	case PubRec:
		commRef := commsem.DataFromRef(rec.CommRef)
		return StepRecDS{
			CommID: commRef.CommID,
			CommRN: commRef.CommRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      PubStep,
			Exp:    poolexp.DataFromExpRec(rec.ValExp),
		}
	case SubRec:
		commRef := commsem.DataFromRef(rec.CommRef)
		return StepRecDS{
			CommID: commRef.CommID,
			CommRN: commRef.CommRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      SubStep,
			Exp:    poolexp.DataFromExpRec(rec.ContExp),
		}
	default:
		panic(ErrRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto StepRecDS) (StepRec, error) {
	switch dto.K {
	case PubStep:
		commRef, err := DataToSemRef(dto)
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
		return PubRec{CommRef: commRef, ChnlID: chnlID, ValExp: valExp}, nil
	case SubStep:
		commRef, err := DataToSemRef(dto)
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
		return SubRec{CommRef: commRef, ChnlID: chnlID, ContExp: contExp}, nil
	default:
		panic(ErrStepKindUnexpected(dto.K))
	}
}
