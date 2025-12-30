package identity

import (
	"github.com/rs/xid"
)

func ConvertToSame(id ADT) ADT {
	return id
}

func ConvertFromString(s string) (ADT, error) {
	xid, err := xid.FromString(s)
	if err != nil {
		return ADT{}, err
	}
	return ADT(xid), nil
}

func ConvertToString(id ADT) string {
	return xid.ID(id).String()
}

func ConvertPtrToStringPtr(id *ADT) *string {
	if id == nil {
		return nil
	}
	s := xid.ID(*id).String()
	return &s
}
