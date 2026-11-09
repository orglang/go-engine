package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	conncore "orglang/go-engine/internal/pool/conn/core"
)

func identityToString(id identity.ADT) string {
	return identity.ConvertToString(id)
}

func dataToConnID(dto conncore.ConnSnapDS) (identity.ADT, error) {
	return identity.ConvertFromString(dto.ConnID)
}

func DataFromRec(rec conncore.ConnRec) conncore.ConnRecDS {
	return conncore.ConnRecDS{
		ConnID:         identity.ConvertToString(rec.ConnID),
		ProviderCompID: identity.ConvertToString(rec.ProviderCompRef.CompID),
		ProviderCompRN: int64(rec.ProviderCompRef.CompRN),
		ClientCompID:   identity.ConvertToString(rec.ClientCompRef.CompID),
		ClientCompRN:   int64(rec.ClientCompRef.CompRN),
		OffsetNr:       int64(rec.OffsetNr),
		ActNr:          int64(rec.ActNr),
	}
}
