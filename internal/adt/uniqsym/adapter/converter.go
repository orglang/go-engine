package adapter

import (
	"database/sql"

	"orglang/go-engine/internal/adt/uniqsym"
)

func ConvertFromNullString(str sql.NullString) (uniqsym.ADT, error) {
	if !str.Valid {
		var empty uniqsym.ADT
		return empty, nil
	}
	return uniqsym.ConvertFromString(str.String)
}

func ConvertToNullString(adt uniqsym.ADT) sql.NullString {
	var empty uniqsym.ADT
	if adt == empty {
		return sql.NullString{}
	}
	return sql.NullString{String: uniqsym.ConvertToString(adt), Valid: true}
}
