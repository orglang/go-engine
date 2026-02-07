package symbol

import (
	"fmt"
)

func ConvertFromString(str string) (ADT, error) {
	if str == "" {
		return "", fmt.Errorf("invalid value: %s", str)
	}
	return ADT(str), nil
}

func ConvertToString(adt ADT) string {
	return string(adt)
}
