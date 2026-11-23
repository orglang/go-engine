package contkey

import (
	jsonv2 "encoding/json/v2"
	"encoding/json/jsontext"
	"hash/fnv"
)

// HashJSON returns the FNV-1a 64-bit content key and canonical JSON bytes
// for value. Callers must include a format version and a type discriminator
// in value when those are relevant to the identity being hashed.
func HashJSON(value any) (ADT, []byte, error) {
	data, err := jsonv2.Marshal(value)
	if err != nil {
		return Zero, nil, err
	}

	canonical := jsontext.Value(data)
	if err := canonical.Canonicalize(jsontext.CanonicalizeRawInts(false)); err != nil {
		return Zero, nil, err
	}

	h := fnv.New64a()
	_, _ = h.Write(canonical)

	return ADT(h.Sum64()), canonical, nil
}
