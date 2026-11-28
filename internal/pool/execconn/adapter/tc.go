package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	execconn "orglang/go-engine/internal/pool/execconn/core"
)

func identityToString(id identity.ADT) string {
	return identity.ConvertToString(id)
}

func dataToConnID(dto ConnSnap) (identity.ADT, error) {
	return identity.ConvertFromString(dto.ConnID)
}

func DataFromRec(rec execconn.ConnRec) ConnRec {
	return ConnRec{
		ConnID:         identity.ConvertToString(rec.ConnID),
		ProviderExecID: identity.ConvertToString(rec.ProviderRootRef.RootID),
		ProviderExecRN: int64(rec.ProviderRootRef.RootRN),
		ClientExecID:   identity.ConvertToString(rec.ClientRootRef.RootID),
		ClientExecRN:   int64(rec.ClientRootRef.RootRN),
		OffsetNr:       int64(rec.OffsetNr),
		ActNr:          int64(rec.ActNr),
	}
}
