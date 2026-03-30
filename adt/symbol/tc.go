package symbol

import (
	"database/sql"
	"fmt"
)

func ConvertFromString(str string) (ADT, error) {
	if str == "" {
		return Unit, fmt.Errorf("invalid value: %s", str)
	}
	return ADT(str), nil
}

func ConvertToString(adt ADT) string {
	return string(adt)
}

func ConvertFromNullString(str sql.NullString) (ADT, error) {
	if str.Valid {
		return ADT(str.String), nil
	}
	return Unit, nil
}

func ConvertToNullString(adt ADT) sql.NullString {
	return sql.NullString{String: string(adt), Valid: true}
}
