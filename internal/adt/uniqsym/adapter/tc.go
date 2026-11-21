package adapter

import (
	"database/sql"

	"orglang/go-engine/internal/adt/uniqsym"
)

func ConvertFromNullString(str sql.NullString) (uniqsym.ADT, error) {
	if !str.Valid {
		return uniqsym.Empty(), nil
	}
	return uniqsym.ConvertFromString(str.String)
}

func ConvertToNullString(adt uniqsym.ADT) sql.NullString {
	if adt.IsEmpty() {
		return sql.NullString{}
	}
	return sql.NullString{String: uniqsym.ConvertToString(adt), Valid: true}
}
