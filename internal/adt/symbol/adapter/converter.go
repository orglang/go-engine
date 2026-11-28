package adapter

import (
	"database/sql"

	"orglang/go-engine/internal/adt/symbol"
)

func ConvertFromNullString(str sql.NullString) (symbol.ADT, error) {
	if str.Valid {
		return symbol.ADT(str.String), nil
	}
	return symbol.Zero, nil
}

func ConvertToNullString(adt symbol.ADT) sql.NullString {
	return sql.NullString{String: string(adt), Valid: true}
}
