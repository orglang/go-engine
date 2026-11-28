package adapter

import (
	"fmt"
	execrootadapter "orglang/go-engine/internal/adt/execroot/adapter"
	"orglang/go-engine/internal/adt/identity"
	execeval "orglang/go-engine/internal/pool/execeval/adapter"

	commact "orglang/go-engine/internal/pool/commact/core"
)

func DataFromStepRec(r commact.ActRec) ActRec {
	switch rec := r.(type) {
	case commact.PubRec:
		rootRef := execrootadapter.DataFromRef(rec.RootRef)
		return ActRec{
			CommID: identity.ConvertToString(rec.CommID),
			ActNr:  int64(rec.ActNr),
			ExecID: rootRef.ExecID,
			ExecRN: rootRef.ExecRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      PubKind,
			Expr:    execeval.DataFromExprRec(rec.ValExp),
		}
	case commact.SubRec:
		rootRef := execrootadapter.DataFromRef(rec.RootRef)
		return ActRec{
			CommID: identity.ConvertToString(rec.CommID),
			ActNr:  int64(rec.ActNr),
			ExecID: rootRef.ExecID,
			ExecRN: rootRef.ExecRN,
			ChnlID: identity.ConvertToString(rec.ChnlID),
			K:      SubKind,
			Expr:    execeval.DataFromExprRec(rec.ContExp),
		}
	default:
		panic(commact.ErrActRecTypeUnexpected(r))
	}
}

func DataToStepRec(dto ActRec) (commact.ActRec, error) {
	switch dto.K {
	case PubKind:
		return DataToPubRec(dto)
	case SubKind:
		return DataToSubRec(dto)
	default:
		panic(fmt.Errorf("act kind unexpected: %v", dto.K))
	}
}
