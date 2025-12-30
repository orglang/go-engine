package qualsym

func ConvertToSame(s ADT) ADT {
	return s
}

func ConvertFromString(s string) (ADT, error) {
	return ADT(s), nil
}

func ConvertToString(s ADT) string {
	return string(s)
}
