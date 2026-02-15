package valkey

func ConvertFromInteger(key int64) (ADT, error) {
	return ADT(key), nil
}

func ConvertToInteger(key ADT) int64 {
	return int64(key)
}
