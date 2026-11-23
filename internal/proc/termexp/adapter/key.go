package adapter

import (
    "orglang/go-engine/internal/adt/contkey"
    termexp "orglang/go-engine/internal/proc/termexp/core"
)

// ContentKey hashes the adapter-owned storage representation of a Definition.
func ContentKey(spec termexp.ExpSpec) (contkey.ADT, []byte, error) {
    value, err := canonicalExpSpec(spec)
    if err != nil {
        return contkey.Zero, nil, err
    }
    return contkey.HashJSON(value)
}
