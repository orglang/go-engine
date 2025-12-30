package revnum

func ConvertToInt(a ADT) int64 {
	return int64(a)
}

func ConvertFromInt(i int64) ADT {
	return ADT(i)
}
