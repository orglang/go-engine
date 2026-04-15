package commturn

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/semcomm"
	"orglang/go-engine/adt/semcomp"
	"orglang/go-engine/pool/termexp"
)

func DataFromStepRec(r TurnRec) TurnRecDS {
	switch rec := r.(type) {
	case PubRec:
		commRef := semcomm.DataFromRef(rec.CommRef)
		implRef := semcomp.DataFromRef(rec.CompRef)
		return TurnRecDS{
			CommID: commRef.CommID,
			CommRN: commRef.CommRN,
			CompID: implRef.CompID,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      PubKind,
			Exp:    termexp.DataFromExpRec(rec.ValExp),
		}
	case SubRec:
		commRef := semcomm.DataFromRef(rec.CommRef)
		implRef := semcomp.DataFromRef(rec.CompRef)
		return TurnRecDS{
			CommID: commRef.CommID,
			CommRN: commRef.CommRN,
			CompID: implRef.CompID,
			ChnlID: identity.ConvertToString(rec.ChnlID),
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
		commRef, err := DataToCommRef(dto)
		if err != nil {
			return nil, err
		}
		implRef, err := DataToCompRef(dto)
		if err != nil {
			return nil, err
		}
		chnlID, err := identity.ConvertFromString(dto.ChnlID)
		if err != nil {
			return nil, err
		}
		valExp, err := termexp.DataToExpRec(dto.Exp)
		if err != nil {
			return nil, err
		}
		return PubRec{CommRef: commRef, CompRef: implRef, ChnlID: chnlID, ValExp: valExp}, nil
	case SubKind:
		commRef, err := DataToCommRef(dto)
		if err != nil {
			return nil, err
		}
		implRef, err := DataToCompRef(dto)
		if err != nil {
			return nil, err
		}
		chnlID, err := identity.ConvertFromString(dto.ChnlID)
		if err != nil {
			return nil, err
		}
		contExp, err := termexp.DataToExpRec(dto.Exp)
		if err != nil {
			return nil, err
		}
		return SubRec{CommRef: commRef, CompRef: implRef, ChnlID: chnlID, ContExp: contExp}, nil
	default:
		panic(ErrStepKindUnexpected(dto.K))
	}
}
