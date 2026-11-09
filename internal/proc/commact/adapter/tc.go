package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"

	commact "orglang/go-engine/internal/proc/commact/core"
	termexp "orglang/go-engine/internal/proc/termexp/adapter"
)

func DataFromStepRec(r commact.ActRec) (commact.ActRecDS, error) {
	switch rec := r.(type) {
	case commact.PubRec:
		exp, err := termexp.DataFromExpRec(rec.ValExp)
		if err != nil {
			return commact.ActRecDS{}, err
		}
		compRef := compsem.DataFromRef(rec.CompRef)
		return commact.ActRecDS{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commact.MsgStep,
			Exp:    exp,
		}, nil
	case commact.SubRec:
		exp, err := termexp.DataFromExpRec(rec.ContExp)
		if err != nil {
			return commact.ActRecDS{}, err
		}
		compRef := compsem.DataFromRef(rec.CompRef)
		return commact.ActRecDS{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			CompID: compRef.CompID,
			CompRN: compRef.CompRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      commact.SvcStep,
			Exp:    exp,
		}, nil
	default:
		panic(commact.ErrActRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto commact.ActRecDS) (commact.ActRec, error) {
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
	case commact.MsgStep:
		return commact.PubRec{ConnID: connID, ActNr: seqnum.ADT(dto.ActNr), CompRef: compRef, ChnlID: chnlID, ValExp: exp}, nil
	case commact.SvcStep:
		return commact.SubRec{ConnID: connID, ActNr: seqnum.ADT(dto.ActNr), CompRef: compRef, ChnlID: chnlID, ContExp: exp}, nil
	default:
		return nil, fmt.Errorf("unexpected Act kind: %v", dto.K)
	}
}
