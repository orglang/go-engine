package valkey

import (
	"database/sql"
)

func ConvertFromInt(key int64) (ADT, error) {
	return ADT(key), nil
}

func ConvertToInt(key ADT) int64 {
	return int64(key)
}

func ConvertToNullInt(a ADT) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(a), Valid: true}
}

func ConvertFromNullInt(n sql.NullInt64) ADT {
	if n.Valid {
		return ADT(n.Int64)
	}
	return Zero
}
