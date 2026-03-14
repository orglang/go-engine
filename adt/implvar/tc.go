package implvar

import (
	"orglang/go-engine/adt/implsem"
)

func ConvertRecToRef(rec VarRec) implsem.SemRef {
	return rec.ImplRef
}

func DataFromVarRec(VarRec2) VarRecDS {
	return VarRecDS{}
}

func DataToVarRec(VarRecDS) (VarRec2, error) {
	return nil, nil
}
