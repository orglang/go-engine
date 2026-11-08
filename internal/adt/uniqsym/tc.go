package uniqsym

import (
	"database/sql"
	"fmt"
	"strings"

	"orglang/go-engine/adt/symbol"
)

const (
	sep = "."
)

func ConvertFromString(str string) (ADT, error) {
	if str == "" {
		return empty, fmt.Errorf("invalid value: %s", str)
	}
	idx := strings.LastIndex(str, sep)
	if idx < 0 {
		return ADT{symbol.New(str), nil}, nil
	}
	sym, err := symbol.ConvertFromString(str[idx+1:])
	if err != nil {
		return empty, err
	}
	ns, err := ConvertFromString(str[:idx])
	if err != nil {
		return empty, err
	}
	return ADT{sym, &ns}, nil
}

func ConvertFromNullString(str sql.NullString) (ADT, error) {
	if !str.Valid {
		return empty, nil
	}
	idx := strings.LastIndex(str.String, sep)
	if idx < 0 {
		return ADT{symbol.New(str.String), nil}, nil
	}
	sym, err := symbol.ConvertFromString(str.String[idx+1:])
	if err != nil {
		return empty, err
	}
	ns, err := ConvertFromString(str.String[:idx])
	if err != nil {
		return empty, err
	}
	return ADT{sym, &ns}, nil
}

func ConvertToString(adt ADT) string {
	if adt == empty {
		panic("invalid value")
	}
	sym := symbol.ConvertToString(adt.sym)
	if adt.ns == nil {
		return sym
	}
	ns := ConvertToString(*adt.ns)
	return ns + sep + sym
}

func ConvertToNullString(adt ADT) sql.NullString {
	if adt == empty {
		return sql.NullString{}
	}
	sym := symbol.ConvertToString(adt.sym)
	if adt.ns == nil {
		return sql.NullString{String: sym, Valid: true}
	}
	ns := ConvertToString(*adt.ns)
	return sql.NullString{String: ns + sep + sym, Valid: true}
}
