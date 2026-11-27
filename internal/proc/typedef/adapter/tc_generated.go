package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

func DefSpecFromGenerated(dto webapi.ProcTypeDefSpec) (typedef.DefSpec, error) {
	qn, err := uniqsym.ConvertFromString(dto.TypeQn)
	if err != nil {
		return typedef.DefSpec{}, err
	}
	exp, err := expSpecFromGenerated(dto.TypeExpr)
	if err != nil {
		return typedef.DefSpec{}, err
	}
	return typedef.DefSpec{TypeQN: qn, TypeExp: exp}, nil
}

func DefSnapToGenerated(snap typedef.DefSnap) webapi.ProcTypeDefSnap {
	return webapi.ProcTypeDefSnap{
		Ref: webapi.AnyTypeRef{
			TypeId: identity.ConvertToString(snap.TypeID),
			DescRn: seqnum.ConvertToInt(snap.DescRef.DescRN),
		},
		Spec: webapi.ProcTypeDefSpec{
			TypeQn:   uniqsym.ConvertToString(snap.DefSpec.TypeQN),
			TypeExpr: expSpecToGenerated(snap.DefSpec.TypeExp),
		},
	}
}

func expSpecFromGenerated(dto webapi.ProcTypeExprSpec) (typedef.ExpSpec, error) {
	kind, err := dto.Discriminator()
	if err != nil {
		return nil, err
	}
	switch kind {
	case "one":
		return typedef.OneSpec{}, nil
	case "link":
		v, err := dto.AsProcTypeLinkSpec()
		if err != nil {
			return nil, err
		}
		qn, err := uniqsym.ConvertFromString(v.Link.TypeQn)
		if err != nil {
			return nil, err
		}
		return typedef.LinkSpec{TypeQN: qn}, nil
	case "tensor":
		v, err := dto.AsProcTypeTensorSpec()
		if err != nil {
			return nil, err
		}
		val, err := expSpecFromGenerated(v.ValExpr)
		if err != nil {
			return nil, err
		}
		cont, err := expSpecFromGenerated(v.ContExpr)
		if err != nil {
			return nil, err
		}
		return typedef.TensorSpec{Val: val, Cont: cont}, nil
	case "lolli":
		v, err := dto.AsProcTypeLolliSpec()
		if err != nil {
			return nil, err
		}
		val, err := expSpecFromGenerated(v.ValExpr)
		if err != nil {
			return nil, err
		}
		cont, err := expSpecFromGenerated(v.ContExpr)
		if err != nil {
			return nil, err
		}
		return typedef.LolliSpec{Val: val, Cont: cont}, nil
	case "plus":
		v, err := dto.AsProcTypePlusSpec()
		if err != nil {
			return nil, err
		}
		choices := make(map[uniqsym.ADT]typedef.ExpSpec, len(v.Choices))
		for _, choice := range v.Choices {
			label, err := uniqsym.ConvertFromString(choice.LabQn)
			if err != nil {
				return nil, err
			}
			cont, err := expSpecFromGenerated(choice.ContExpr)
			if err != nil {
				return nil, err
			}
			choices[label] = cont
		}
		return typedef.PlusSpec{Choices: choices}, nil
	case "with":
		v, err := dto.AsProcTypeWithSpec()
		if err != nil {
			return nil, err
		}
		choices := make(map[uniqsym.ADT]typedef.ExpSpec, len(v.Choices))
		for _, choice := range v.Choices {
			label, err := uniqsym.ConvertFromString(choice.LabQn)
			if err != nil {
				return nil, err
			}
			cont, err := expSpecFromGenerated(choice.ContExpr)
			if err != nil {
				return nil, err
			}
			choices[label] = cont
		}
		return typedef.WithSpec{Choices: choices}, nil
	case "up":
		v, err := dto.AsProcTypeUpSpec()
		if err != nil {
			return nil, err
		}
		cont, err := expSpecFromGenerated(v.ContExpr)
		if err != nil {
			return nil, err
		}
		return typedef.UpSpec{Cont: cont}, nil
	case "down":
		v, err := dto.AsProcTypeDownSpec()
		if err != nil {
			return nil, err
		}
		cont, err := expSpecFromGenerated(v.ContExpr)
		if err != nil {
			return nil, err
		}
		return typedef.DownSpec{Cont: cont}, nil
	default:
		return nil, fmt.Errorf("unsupported process type expression kind %q", kind)
	}
}

func expSpecToGenerated(spec typedef.ExpSpec) webapi.ProcTypeExprSpec {
	var dto webapi.ProcTypeExprSpec
	switch spec := spec.(type) {
	case typedef.OneSpec:
		_ = dto.FromProcTypeOneSpec(webapi.ProcTypeOneSpec{})
	case typedef.LinkSpec:
		_ = dto.FromProcTypeLinkSpec(webapi.ProcTypeLinkSpec{
			Link: webapi.AnyTypeLinkSpec{TypeQn: uniqsym.ConvertToString(spec.TypeQN)},
		})
	case typedef.TensorSpec:
		_ = dto.FromProcTypeTensorSpec(webapi.ProcTypeTensorSpec{
			ValExpr:  expSpecToGenerated(spec.Val),
			ContExpr: expSpecToGenerated(spec.Cont),
		})
	case typedef.LolliSpec:
		_ = dto.FromProcTypeLolliSpec(webapi.ProcTypeLolliSpec{
			ValExpr:  expSpecToGenerated(spec.Val),
			ContExpr: expSpecToGenerated(spec.Cont),
		})
	case typedef.PlusSpec:
		choices := make([]webapi.ProcTypeChoiceSpec, 0, len(spec.Choices))
		for label, cont := range spec.Choices {
			choices = append(choices, webapi.ProcTypeChoiceSpec{
				LabQn:    uniqsym.ConvertToString(label),
				ContExpr: expSpecToGenerated(cont),
			})
		}
		_ = dto.FromProcTypePlusSpec(webapi.ProcTypePlusSpec{Choices: choices})
	case typedef.WithSpec:
		choices := make([]webapi.ProcTypeChoiceSpec, 0, len(spec.Choices))
		for label, cont := range spec.Choices {
			choices = append(choices, webapi.ProcTypeChoiceSpec{
				LabQn:    uniqsym.ConvertToString(label),
				ContExpr: expSpecToGenerated(cont),
			})
		}
		_ = dto.FromProcTypeWithSpec(webapi.ProcTypeWithSpec{Choices: choices})
	case typedef.UpSpec:
		_ = dto.FromProcTypeUpSpec(webapi.ProcTypeUpSpec{
			ContExpr: expSpecToGenerated(spec.Cont),
		})
	case typedef.DownSpec:
		_ = dto.FromProcTypeDownSpec(webapi.ProcTypeDownSpec{
			ContExpr: expSpecToGenerated(spec.Cont),
		})
	default:
		panic(typedef.ErrSpecTypeUnexpected(spec))
	}
	return dto
}
