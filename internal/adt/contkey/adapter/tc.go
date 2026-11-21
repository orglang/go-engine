package adapter

import (
	"database/sql"

	"orglang/go-engine/internal/adt/contkey"
)

func ConvertToNullInt(a contkey.ADT) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(a), Valid: true}
}

func ConvertFromNullInt(n sql.NullInt64) contkey.ADT {
	if n.Valid {
		return contkey.ADT(n.Int64)
	}
	return contkey.Zero
}
