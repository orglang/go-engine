package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/execroot"
	execrootadapter "orglang/go-engine/internal/adt/execroot/adapter"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"

	commact "orglang/go-engine/internal/proc/commact/core"
	execeval "orglang/go-engine/internal/proc/execeval/adapter"
)

func DataFromStepRec(r commact.ActRec) (ActRec, error) {
	switch rec := r.(type) {
	case commact.PubRec:
		exp, err := execeval.DataFromExpRec(rec.ValExp)
		if err != nil {
			return ActRec{}, err
		}
		rootRef := execrootadapter.DataFromRef(rec.RootRef)
		return ActRec{
			CommID: identity.ConvertToString(rec.CommID),
			ActNr:  int64(rec.ActNr),
			ExecID: rootRef.ExecID,
			ExecRN: rootRef.ExecRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      MsgStep,
			Exp:    exp,
		}, nil
	case commact.SubRec:
		exp, err := execeval.DataFromExpRec(rec.ContExp)
		if err != nil {
			return ActRec{}, err
		}
		rootRef := execrootadapter.DataFromRef(rec.RootRef)
		return ActRec{
			CommID: identity.ConvertToString(rec.CommID),
			ActNr:  int64(rec.ActNr),
			ExecID: rootRef.ExecID,
			ExecRN: rootRef.ExecRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      SvcStep,
			Exp:    exp,
		}, nil
	default:
		panic(commact.ErrActRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto ActRec) (commact.ActRec, error) {
	commID, err := identity.ConvertFromString(dto.CommID)
	if err != nil {
		return nil, err
	}
	rootID, err := identity.ConvertFromString(dto.ExecID)
	if err != nil {
		return nil, err
	}
	chnlID, err := identity.ConvertFromString(dto.ChnlID)
	if err != nil {
		return nil, err
	}
	exp, err := execeval.DataToExpRec(dto.Exp)
	if err != nil {
		return nil, err
	}
	rootRef := execroot.RootRef{ExecID: rootID, ExecRN: seqnum.ADT(dto.ExecRN)}
	switch dto.K {
	case MsgStep:
		return commact.PubRec{CommID: commID, ActNr: seqnum.ADT(dto.ActNr), RootRef: rootRef, ChnlID: chnlID, ValExp: exp}, nil
	case SvcStep:
		return commact.SubRec{CommID: commID, ActNr: seqnum.ADT(dto.ActNr), RootRef: rootRef, ChnlID: chnlID, ContExp: exp}, nil
	default:
		return nil, fmt.Errorf("unexpected Act kind: %v", dto.K)
	}
}
