package uniqsym

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"orglang/go-engine/internal/adt/symbol"
)

const (
	sep = "."
	maxLength = 512
)

var qualifiedSymbolPattern = regexp.MustCompile(`^[0-9A-Za-z_-]+(\.[0-9A-Za-z_-]+)*$`)

// ConvertFromString constructs a qualified symbol and rejects values outside
// the domain's qualified-symbol syntax.
func ConvertFromString(str string) (ADT, error) {
	if len(str) == 0 || len(str) > maxLength || !qualifiedSymbolPattern.MatchString(str) {
		return empty, fmt.Errorf("invalid qualified symbol: %q", str)
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
	return ConvertFromString(str.String)
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
