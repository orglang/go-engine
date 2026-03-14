package seqnum

type ADT int64

func New() ADT {
	return ADT(1)
}

func Next(sn ADT) ADT {
	return sn + 1
}

func (sn ADT) Next() ADT {
	return sn + 1
}
