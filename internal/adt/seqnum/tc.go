package seqnum

func ConvertToInt(a ADT) int64 {
	return int64(a)
}

func ConvertFromInt(n int64) ADT {
	return ADT(n)
}
