package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	execcomm "orglang/go-engine/internal/pool/execcomm/core"
)

func identityToString(id identity.ADT) string {
	return identity.ConvertToString(id)
}

func dataToCommID(dto CommSnap) (identity.ADT, error) {
	return identity.ConvertFromString(dto.CommID)
}

func DataFromRec(rec execcomm.CommRec) CommRec {
	return CommRec{
		CommID:         identity.ConvertToString(rec.CommID),
		ProviderExecID: identity.ConvertToString(rec.ProviderRootRef.ExecID),
		ProviderExecRN: int64(rec.ProviderRootRef.ExecRN),
		ClientExecID:   identity.ConvertToString(rec.ClientRootRef.ExecID),
		ClientExecRN:   int64(rec.ClientRootRef.ExecRN),
		OffsetNr:       int64(rec.OffsetNr),
		ActNr:          int64(rec.ActNr),
	}
}
