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
		ProviderExecID: identity.ConvertToString(rec.ProviderExecRef.CompID),
		ProviderExecRN: int64(rec.ProviderExecRef.CompRN),
		ClientExecID:   identity.ConvertToString(rec.ClientExecRef.CompID),
		ClientExecRN:   int64(rec.ClientExecRef.CompRN),
		OffsetNr:       int64(rec.OffsetNr),
		ActNr:          int64(rec.ActNr),
	}
}
