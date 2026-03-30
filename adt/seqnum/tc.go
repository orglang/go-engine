package seqnum

import (
	"database/sql"
)

func ConvertToInt(a ADT) int64 {
	return int64(a)
}

func ConvertFromInt(n int64) ADT {
	return ADT(n)
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
