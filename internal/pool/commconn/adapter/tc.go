package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	commconn "orglang/go-engine/internal/pool/commconn/core"
)

func identityToString(id identity.ADT) string {
	return identity.ConvertToString(id)
}

func dataToConnID(dto commconn.ConnSnapDS) (identity.ADT, error) {
	return identity.ConvertFromString(dto.ConnID)
}
