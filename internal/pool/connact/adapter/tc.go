package adapter

import (
	"fmt"
	execrootadapter "orglang/go-engine/internal/adt/execroot/adapter"
	"orglang/go-engine/internal/adt/identity"
	execeval "orglang/go-engine/internal/pool/execeval/adapter"

	connact "orglang/go-engine/internal/pool/connact/core"
)

func DataFromStepRec(r connact.ActRec) ActRec {
	switch rec := r.(type) {
	case connact.PubRec:
		rootRef := execrootadapter.DataFromRef(rec.RootRef)
		return ActRec{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			ExecID: rootRef.RootID,
			ExecRN: rootRef.RootRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      PubKind,
			Exp:    execeval.DataFromExpRec(rec.ValExp),
		}
	case connact.SubRec:
		rootRef := execrootadapter.DataFromRef(rec.RootRef)
		return ActRec{
			ConnID: identity.ConvertToString(rec.ConnID),
			ActNr:  int64(rec.ActNr),
			ExecID: rootRef.RootID,
			ExecRN: rootRef.RootRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      SubKind,
			Exp:    execeval.DataFromExpRec(rec.ContExp),
		}
	default:
		panic(connact.ErrActRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto ActRec) (connact.ActRec, error) {
	switch dto.K {
	case PubKind:
		return DataToPubRec(dto)
	case SubKind:
		return DataToSubRec(dto)
	default:
		panic(fmt.Errorf("act kind unexpected: %v", dto.K))
	}
}
