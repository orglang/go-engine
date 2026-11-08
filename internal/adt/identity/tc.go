package identity

import (
	"database/sql"
	"orglang/go-engine/internal/adt/option"

	"github.com/rs/xid"
)

func ConvertFromString(s string) (ADT, error) {
	id, err := xid.FromString(s)
	if err != nil {
		return ADT{}, err
	}
	return ADT(id), nil
}

func ConvertToString(id ADT) string {
	return xid.ID(id).String()
}

func ConvertPtrToStringPtr(id *ADT) *string {
	if id == nil {
		return nil
	}
	s := xid.ID(*id).String()
	return &s
}

func ConvertFromNullString(s sql.NullString) (ADT, error) {
	if !s.Valid || s.String == "" {
		return ADT{}, nil
	}
	return ConvertFromString(s.String)
}

func ConvertToNullString(id ADT) sql.NullString {
	if id.IsEmpty() {
		return sql.NullString{}
	}
	return sql.NullString{String: ConvertToString(id), Valid: true}
}

func ConvertOptionToNullString(id option.ADT[ADT]) sql.Null[string] {
	if id.IsEmpty() {
		return sql.Null[string]{}
	}
	return sql.Null[string]{V: ConvertToString(id.Get()), Valid: true}
}
