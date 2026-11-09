package adapter

import (
	sdktermvar "github.com/orglang/go-sdk/adt/termvar"
	sdktypesem "github.com/orglang/go-sdk/adt/typesem"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/typesem"
)

func MsgFromTermVarRec(rec termvar.VarRec) sdktermvar.VarRec {
	return sdktermvar.VarRec{
		TypeRef: sdktypesem.SemRef{
			TypeID: identity.ConvertToString(rec.TypeRef.TypeID),
			TypeRN: seqnum.ConvertToInt(rec.TypeRef.TypeRN),
		},
		ChnlPH: symbol.ConvertToString(rec.ChnlPH),
		ExpVK:  contkey.ConvertToInt(rec.ExpCK),
	}
}

func MsgToTermVarRec(dto sdktermvar.VarRec) (termvar.VarRec, error) {
	typeID, err := identity.ConvertFromString(dto.TypeRef.TypeID)
	if err != nil {
		return termvar.VarRec{}, err
	}
	typeRN := seqnum.ConvertFromInt(dto.TypeRef.TypeRN)
	chnlPH, err := symbol.ConvertFromString(dto.ChnlPH)
	if err != nil {
		return termvar.VarRec{}, err
	}
	expCK, err := contkey.ConvertFromInt(dto.ExpVK)
	if err != nil {
		return termvar.VarRec{}, err
	}

	return termvar.VarRec{
		TypeRef: typesem.SemRef{TypeID: typeID, TypeRN: typeRN},
		ChnlPH:  chnlPH,
		ExpCK:   expCK,
	}, nil
}
