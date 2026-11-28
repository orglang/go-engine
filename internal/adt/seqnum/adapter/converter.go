package adapter

import (
	"database/sql"

	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"
)

func ConvertToNullInt(a seqnum.ADT) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(a), Valid: true}
}

func ConvertOptionToNullInt(a option.ADT[seqnum.ADT]) sql.Null[int64] {
	if a == nil || a.IsEmpty() {
		return sql.Null[int64]{}
	}
	return sql.Null[int64]{V: seqnum.ConvertToInt(a.Get()), Valid: true}
}

func ConvertFromNullInt(n sql.NullInt64) seqnum.ADT {
	if n.Valid {
		return seqnum.ADT(n.Int64)
	}
	return seqnum.Zero
}
