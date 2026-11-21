package contkey

func ConvertFromInt(key int64) (ADT, error) {
	return ADT(key), nil
}

func ConvertToInt(key ADT) int64 {
	return int64(key)
}
