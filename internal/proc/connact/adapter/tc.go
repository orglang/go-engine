package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/execroot"
	execrootadapter "orglang/go-engine/internal/adt/execroot/adapter"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"

	connact "orglang/go-engine/internal/proc/connact/core"
	execeval "orglang/go-engine/internal/proc/execeval/adapter"
)

func DataFromStepRec(r connact.ActRec) (ActRec, error) {
	switch rec := r.(type) {
	case connact.PubRec:
		exp, err := execeval.DataFromExpRec(rec.ValExp)
		if err != nil {
			return ActRec{}, err
		}
		rootRef := execrootadapter.DataFromRef(rec.RootRef)
		return ActRec{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			ExecID: rootRef.RootID,
			ExecRN: rootRef.RootRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      MsgStep,
			Exp:    exp,
		}, nil
	case connact.SubRec:
		exp, err := execeval.DataFromExpRec(rec.ContExp)
		if err != nil {
			return ActRec{}, err
		}
		rootRef := execrootadapter.DataFromRef(rec.RootRef)
		return ActRec{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			ExecID: rootRef.RootID,
			ExecRN: rootRef.RootRN,
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
	compID, err := identity.ConvertFromString(dto.ExecID)
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
	rootRef := execroot.RootRef{RootID: compID, RootRN: seqnum.ADT(dto.ExecRN)}
	switch dto.K {
	case MsgStep:
		return connact.PubRec{ConnID: connID, ActNr: seqnum.ADT(dto.ActNr), RootRef: rootRef, ChnlID: chnlID, ValExp: exp}, nil
	case SvcStep:
		return connact.SubRec{ConnID: connID, ActNr: seqnum.ADT(dto.ActNr), RootRef: rootRef, ChnlID: chnlID, ContExp: exp}, nil
	default:
		return nil, fmt.Errorf("unexpected Act kind: %v", dto.K)
	}
}
