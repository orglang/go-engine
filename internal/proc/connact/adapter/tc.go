package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/comproot"
	compsemadapter "orglang/go-engine/internal/adt/comproot/adapter"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"

	connact "orglang/go-engine/internal/proc/connact/core"
	compexec "orglang/go-engine/internal/proc/compexec/adapter"
)

func DataFromStepRec(r connact.ActRec) (ActRec, error) {
	switch rec := r.(type) {
	case connact.PubRec:
		exp, err := compexec.DataFromExpRec(rec.ValExp)
		if err != nil {
			return ActRec{}, err
		}
		compRef := compsemadapter.DataFromRef(rec.CompRef)
		return ActRec{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      MsgStep,
			Exp:    exp,
		}, nil
	case connact.SubRec:
		exp, err := compexec.DataFromExpRec(rec.ContExp)
		if err != nil {
			return ActRec{}, err
		}
		compRef := compsemadapter.DataFromRef(rec.CompRef)
		return ActRec{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      SvcStep,
			Exp:    exp,
		}, nil
	default:
		panic(connact.ErrActRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto ActRec) (connact.ActRec, error) {
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
	exp, err := compexec.DataToExpRec(dto.Exp)
	if err != nil {
		return nil, err
	}
	compRef := comproot.SemRef{CompID: compID, CompRN: seqnum.ADT(dto.CompRN)}
	switch dto.K {
	case MsgStep:
		return connact.PubRec{ConnID: connID, ActNr: seqnum.ADT(dto.ActNr), CompRef: compRef, ChnlID: chnlID, ValExp: exp}, nil
	case SvcStep:
		return connact.SubRec{ConnID: connID, ActNr: seqnum.ADT(dto.ActNr), CompRef: compRef, ChnlID: chnlID, ContExp: exp}, nil
	default:
		return nil, fmt.Errorf("unexpected Act kind: %v", dto.K)
	}
}
