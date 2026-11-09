package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	compconncore "orglang/go-engine/internal/pool/compconn/core"
)

func identityToString(id identity.ADT) string {
	return identity.ConvertToString(id)
}

func dataToConnID(dto compconncore.ConnSnapDS) (identity.ADT, error) {
	return identity.ConvertFromString(dto.ConnID)
}

func DataFromRec(rec compconncore.ConnRec) compconncore.ConnRecDS {
	return compconncore.ConnRecDS{
		ConnID:         identity.ConvertToString(rec.ConnID),
		ProviderCompID: identity.ConvertToString(rec.ProviderCompRef.CompID),
		ProviderCompRN: int64(rec.ProviderCompRef.CompRN),
		ClientCompID:   identity.ConvertToString(rec.ClientCompRef.CompID),
		ClientCompRN:   int64(rec.ClientCompRef.CompRN),
		OffsetNr:       int64(rec.OffsetNr),
		ActNr:          int64(rec.ActNr),
	}
}
