package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/connsem"
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
		connRef := connsem.DataFromRef(rec.ConnRef)
		compRef := compsem.DataFromRef(rec.CompRef)
		return commturn.TurnRecDS{
			ConnID: connRef.ConnID,
			ConnRN: connRef.ConnRN,
			CompID: compRef.CompID,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commturn.MsgStep,
			Exp:    exp,
		}, nil
	case commturn.SubRec:
		exp, err := termexp.DataFromExpRec(rec.ContExp)
		if err != nil {
			return commturn.TurnRecDS{}, err
		}
		connRef := connsem.DataFromRef(rec.ConnRef)
		compRef := compsem.DataFromRef(rec.CompRef)
		return commturn.TurnRecDS{
			ConnID: connRef.ConnID,
			ConnRN: connRef.ConnRN,
			CompID: compRef.CompID,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commturn.SvcStep,
			Exp:    exp,
		}, nil
	default:
		panic(commturn.ErrRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto commturn.TurnRecDS) (commturn.TurnRec, error) {
	commID, err := identity.ConvertFromString(dto.ConnID)
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
	connRef := connsem.SemRef{ConnID: commID, ConnRN: seqnum.ADT(dto.ConnRN)}
	compRef := compsem.SemRef{CompID: compID}
	switch dto.K {
	case commturn.MsgStep:
		return commturn.PubRec{ConnRef: connRef, CompRef: compRef, ChnlID: chnlID, ValExp: exp}, nil
	case commturn.SvcStep:
		return commturn.SubRec{ConnRef: connRef, CompRef: compRef, ChnlID: chnlID, ContExp: exp}, nil
	default:
		return nil, fmt.Errorf("unexpected Act kind: %v", dto.K)
	}
}
