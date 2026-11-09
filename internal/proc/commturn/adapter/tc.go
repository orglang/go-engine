package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"

	commturn "orglang/go-engine/internal/proc/commturn/core"
	termexp "orglang/go-engine/internal/proc/termexp/adapter"
)

func DataFromStepRec(r commturn.TurnRec) (commturn.TurnRecDS, error) {
	switch rec := r.(type) {
	case commturn.PubRec:
		exp, err := termexp.DataFromExpRec(rec.ValExp)
		if err != nil {
			return commturn.TurnRecDS{}, err
		}
		compRef := compsem.DataFromRef(rec.CompRef)
		return commturn.TurnRecDS{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commturn.MsgStep,
			Exp:    exp,
		}, nil
	case commturn.SubRec:
		exp, err := termexp.DataFromExpRec(rec.ContExp)
		if err != nil {
			return commturn.TurnRecDS{}, err
		}
		compRef := compsem.DataFromRef(rec.CompRef)
		return commturn.TurnRecDS{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commturn.SvcStep,
			Exp:    exp,
		}, nil
	default:
		panic(commturn.ErrRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto commturn.TurnRecDS) (commturn.TurnRec, error) {
	connID, err := identity.ConvertFromString(dto.ConnID)
	if err != nil {
		return nil, err
	}
	compID, err := identity.ConvertFromString(dto.CompID)
	if err != nil {
		return nil, err
	}
	chnlID, err := identity.ConvertFromString(dto.ChnlID)
	if err != nil {
		return nil, err
	}
	exp, err := termexp.DataToExpRec(dto.Exp)
	if err != nil {
		return nil, err
	}
	compRef := compsem.SemRef{CompID: compID, CompRN: seqnum.ADT(dto.CompRN)}
	switch dto.K {
	case commturn.MsgStep:
		return commturn.PubRec{ConnID: connID, ActNr: seqnum.ADT(dto.ActNr), CompRef: compRef, ChnlID: chnlID, ValExp: exp}, nil
	case commturn.SvcStep:
		return commturn.SubRec{ConnID: connID, ActNr: seqnum.ADT(dto.ActNr), CompRef: compRef, ChnlID: chnlID, ContExp: exp}, nil
	default:
		return nil, fmt.Errorf("unexpected Act kind: %v", dto.K)
	}
}
