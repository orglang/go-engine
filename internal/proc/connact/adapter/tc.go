package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"

	connact "orglang/go-engine/internal/proc/connact/core"
	termexp "orglang/go-engine/internal/proc/termexp/adapter"
)

func DataFromStepRec(r connact.ActRec) (connact.ActRecDS, error) {
	switch rec := r.(type) {
	case connact.PubRec:
		exp, err := termexp.DataFromExpRec(rec.ValExp)
		if err != nil {
			return connact.ActRecDS{}, err
		}
		compRef := compsem.DataFromRef(rec.CompRef)
		return connact.ActRecDS{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      connact.MsgStep,
			Exp:    exp,
		}, nil
	case connact.SubRec:
		exp, err := termexp.DataFromExpRec(rec.ContExp)
		if err != nil {
			return connact.ActRecDS{}, err
		}
		compRef := compsem.DataFromRef(rec.CompRef)
		return connact.ActRecDS{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      connact.SvcStep,
			Exp:    exp,
		}, nil
	default:
		panic(connact.ErrActRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto connact.ActRecDS) (connact.ActRec, error) {
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
	case connact.MsgStep:
		return connact.PubRec{ConnID: connID, ActNr: seqnum.ADT(dto.ActNr), CompRef: compRef, ChnlID: chnlID, ValExp: exp}, nil
	case connact.SvcStep:
		return connact.SubRec{ConnID: connID, ActNr: seqnum.ADT(dto.ActNr), CompRef: compRef, ChnlID: chnlID, ContExp: exp}, nil
	default:
		return nil, fmt.Errorf("unexpected Act kind: %v", dto.K)
	}
}
