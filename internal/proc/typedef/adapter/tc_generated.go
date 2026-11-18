package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	typedef "orglang/go-engine/internal/proc/typedef/core"
	typeexp "orglang/go-engine/internal/proc/typeexp/core"
)

func DefSpecFromGenerated(dto webapi.ProcTypeDefSpec) (typedef.DefSpec, error) {
	qn, err := uniqsym.ConvertFromString(dto.TypeQn)
	if err != nil { return typedef.DefSpec{}, err }
	exp, err := expSpecFromGenerated(dto.TypeExpr)
	if err != nil { return typedef.DefSpec{}, err }
	return typedef.DefSpec{TypeQN: qn, TypeExp: exp}, nil
}

func DefSnapToGenerated(snap typedef.DefSnap) webapi.ProcTypeDefSnap {
	return webapi.ProcTypeDefSnap{
		Ref: webapi.AnyTypeRef{TypeId: identity.ConvertToString(snap.TypeRef.TypeID), TypeRn: seqnum.ConvertToInt(snap.TypeRef.TypeRN)},
		Spec: webapi.ProcTypeDefSpec{TypeQn: uniqsym.ConvertToString(snap.DefSpec.TypeQN), TypeExpr: expSpecToGenerated(snap.DefSpec.TypeExp)},
	}
}

func expSpecFromGenerated(dto webapi.ProcTypeExprSpec) (typeexp.ExpSpec, error) {
	kind, err := dto.Discriminator()
	if err != nil { return nil, err }
	switch kind {
	case "one":
		return typeexp.OneSpec{}, nil
	case "link":
		v, err := dto.AsProcTypeLinkSpec(); if err != nil { return nil, err }
		qn, err := uniqsym.ConvertFromString(v.Link.TypeQn); if err != nil { return nil, err }
		return typeexp.LinkSpec{TypeQN: qn}, nil
	case "tensor":
		v, err := dto.AsProcTypeTensorSpec(); if err != nil { return nil, err }
		val, err := expSpecFromGenerated(v.ValExpr); if err != nil { return nil, err }
		cont, err := expSpecFromGenerated(v.ContExpr); if err != nil { return nil, err }
		return typeexp.TensorSpec{Val: val, Cont: cont}, nil
	case "lolli":
		v, err := dto.AsProcTypeLolliSpec(); if err != nil { return nil, err }
		val, err := expSpecFromGenerated(v.ValExpr); if err != nil { return nil, err }
		cont, err := expSpecFromGenerated(v.ContExpr); if err != nil { return nil, err }
		return typeexp.LolliSpec{Val: val, Cont: cont}, nil
	case "plus":
		v, err := dto.AsProcTypePlusSpec(); if err != nil { return nil, err }
		choices := make(map[uniqsym.ADT]typeexp.ExpSpec, len(v.Choices))
		for _, choice := range v.Choices {
			label, err := uniqsym.ConvertFromString(choice.LabQn); if err != nil { return nil, err }
			cont, err := expSpecFromGenerated(choice.ContExpr); if err != nil { return nil, err }
			choices[label] = cont
		}
		return typeexp.PlusSpec{Choices: choices}, nil
	case "with":
		v, err := dto.AsProcTypeWithSpec(); if err != nil { return nil, err }
		choices := make(map[uniqsym.ADT]typeexp.ExpSpec, len(v.Choices))
		for _, choice := range v.Choices {
			label, err := uniqsym.ConvertFromString(choice.LabQn); if err != nil { return nil, err }
			cont, err := expSpecFromGenerated(choice.ContExpr); if err != nil { return nil, err }
			choices[label] = cont
		}
		return typeexp.WithSpec{Choices: choices}, nil
	default:
		return nil, fmt.Errorf("unsupported process type expression kind %q", kind)
	}
}

func expSpecToGenerated(spec typeexp.ExpSpec) webapi.ProcTypeExprSpec {
	var dto webapi.ProcTypeExprSpec
	switch spec := spec.(type) {
	case typeexp.OneSpec:
		_ = dto.FromProcTypeOneSpec(webapi.ProcTypeOneSpec{})
	case typeexp.LinkSpec:
		_ = dto.FromProcTypeLinkSpec(webapi.ProcTypeLinkSpec{Link: webapi.AnyTypeLinkSpec{TypeQn: uniqsym.ConvertToString(spec.TypeQN)}})
	case typeexp.TensorSpec:
		_ = dto.FromProcTypeTensorSpec(webapi.ProcTypeTensorSpec{ValExpr: expSpecToGenerated(spec.Val), ContExpr: expSpecToGenerated(spec.Cont)})
	case typeexp.LolliSpec:
		_ = dto.FromProcTypeLolliSpec(webapi.ProcTypeLolliSpec{ValExpr: expSpecToGenerated(spec.Val), ContExpr: expSpecToGenerated(spec.Cont)})
	case typeexp.PlusSpec:
		choices := make([]webapi.ProcTypeChoiceSpec, 0, len(spec.Choices))
		for label, cont := range spec.Choices { choices = append(choices, webapi.ProcTypeChoiceSpec{LabQn: uniqsym.ConvertToString(label), ContExpr: expSpecToGenerated(cont)}) }
		_ = dto.FromProcTypePlusSpec(webapi.ProcTypePlusSpec{Choices: choices})
	case typeexp.WithSpec:
		choices := make([]webapi.ProcTypeChoiceSpec, 0, len(spec.Choices))
		for label, cont := range spec.Choices { choices = append(choices, webapi.ProcTypeChoiceSpec{LabQn: uniqsym.ConvertToString(label), ContExpr: expSpecToGenerated(cont)}) }
		_ = dto.FromProcTypeWithSpec(webapi.ProcTypeWithSpec{Choices: choices})
	case typeexp.UpSpec:
		_ = dto.FromProcTypeUpSpec(webapi.ProcTypeUpSpec{ContExpr: expSpecToGenerated(spec.Cont)})
	case typeexp.DownSpec:
		_ = dto.FromProcTypeDownSpec(webapi.ProcTypeDownSpec{ContExpr: expSpecToGenerated(spec.Cont)})
	default:
		panic(typeexp.ErrSpecTypeUnexpected(spec))
	}
	return dto
}
