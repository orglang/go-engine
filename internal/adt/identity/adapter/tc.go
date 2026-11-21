package adapter

import (
	"database/sql"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
)

func ConvertToNullString(id identity.ADT) sql.NullString {
	if id.IsEmpty() {
		return sql.NullString{}
	}
	return sql.NullString{String: id.String(), Valid: true}
}

func ConvertFromNullString(dto sql.NullString) (identity.ADT, error) {
	if dto.Valid {
		return identity.ConvertFromString(dto.String)
	}
	return identity.Empty(), nil
}

func ConvertPtrToNullString(id *identity.ADT) sql.NullString {
	if id == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: id.String(), Valid: true}
}

func ConvertOptionToNullString(id option.ADT[identity.ADT]) sql.Null[string] {
	if id.IsEmpty() {
		return sql.Null[string]{}
	}
	return sql.Null[string]{V: identity.ConvertToString(id.Get()), Valid: true}
}
