package core

import typedef "orglang/go-engine/internal/pool/typedef/core"

func ConvertSpecToRec(s ExpSpec) (ExpRec, error) { return typedef.ConvertSpecToRec(s) }
func ConvertRecToSpec(r ExpRec) ExpSpec { return typedef.ConvertRecToSpec(r) }
