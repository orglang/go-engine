package identity

import (
	"github.com/rs/xid"
)

func ConvertFromString(s string) (ADT, error) {
	id, err := xid.FromString(s)
	if err != nil {
		return ADT{}, err
	}
	return ADT(id), nil
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
