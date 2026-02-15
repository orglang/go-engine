package uniqsym

import (
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
