package adapter

import (
	"fmt"
	"slices"
	"strings"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	commact "orglang/go-engine/internal/proc/commact/core"
)

func DataFromExprRec(r commact.ExprRec) (ExprRec, error) {
	switch rec := r.(type) {
	case commact.CloseRec:
		return ExprRec{
			K:     CloseExpr,
			Close: &CloseRec{X: symbol.ConvertToString(rec.ContChnlPH)},
		}, nil
	case commact.WaitRec:
		dto, err := dataFromExprSpec(rec.ContExpr)
		if err != nil {
			return ExprRec{}, err
		}
		return ExprRec{
			K: WaitExpr,
			Wait: &WaitRec{
				X:      symbol.ConvertToString(rec.ContChnlPH),
				ContES: dto,
			},
		}, nil
	case commact.SendRec:
		commID := identity.ConvertToString(rec.CommID)
		return ExprRec{
			K: SendExpr,
			Send: &SendRec{
				X:         symbol.ConvertToString(rec.CommChnlPH),
				A:         identity.ConvertToString(rec.ContChnlID),
				B:         identity.ConvertToString(rec.ValChnlID),
				CommID:    commID,
				ValExprCK: int64(rec.ValExprCK),
			},
		}, nil
	case commact.RecvRec:
		dto, err := dataFromExprSpec(rec.ContExpr)
		if err != nil {
			return ExprRec{}, err
		}
		return ExprRec{
			K: RecvExpr,
			Recv: &RecvRec{
				X:      symbol.ConvertToString(rec.CommChnlPH),
				A:      identity.ConvertToString(rec.ContChnlID),
				Y:      symbol.ConvertToString(rec.NewChnlPH),
				ContES: dto,
			},
		}, nil
	case commact.LabRec:
		return ExprRec{
			K: LabExpr,
			Lab: &LabRec{
				X:         symbol.ConvertToString(rec.CommChnlPH),
				A:         identity.ConvertToString(rec.ContChnlID),
				Label:     uniqsym.ConvertToString(rec.ValLabQN),
				ValExprCK: int64(rec.ValExprCK),
			},
		}, nil
	case commact.CaseRec:
		brs := []BranchRec{}
		for l, cont := range rec.ContExprs {
			dto, err := dataFromExprSpec(cont)
			if err != nil {
				return ExprRec{}, err
			}
			brs = append(brs, BranchRec{Label: uniqsym.ConvertToString(l), ContES: dto})
		}
		return ExprRec{
			K: CaseExpr,
			Case: &CaseRec{
				X:        symbol.ConvertToString(rec.CommChnlPH),
				A:        identity.ConvertToString(rec.ContChnlID),
				Branches: brs,
			},
		}, nil
	case commact.FwdRec:
		return ExprRec{
			K: FwdExpr,
			Fwd: &FwdRec{
				X: symbol.ConvertToString(rec.CommChnlPH),
				B: identity.ConvertToString(rec.ContChnlID),
			},
		}, nil
	default:
		panic(commact.ErrExprTypeUnexpected(rec))
	}
}

func DataToExprRec(dto ExprRec) (commact.ExprRec, error) {
	switch dto.K {
	case CloseExpr:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return commact.CloseRec{ContChnlPH: a}, nil
	case WaitExpr:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExprSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return commact.WaitRec{ContChnlPH: x, ContExpr: cont}, nil
	case SendExpr:
		x, err := symbol.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		a, err := identity.ConvertFromString(dto.Send.A)
		if err != nil {
			return nil, err
		}
		valID, err := identity.ConvertFromString(dto.Send.B)
		if err != nil {
			return nil, err
		}
		commID, err := identity.ConvertFromString(dto.Send.CommID)
		if err != nil {
			return nil, err
		}
		return commact.SendRec{
			CommChnlPH: x,
			ContChnlID: a,
			ValChnlID:  valID,
			ValExprCK:  contkey.ADT(dto.Send.ValExprCK),
			CommID:     commID,
		}, nil
	case RecvExpr:
		x, err := symbol.ConvertFromString(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		a, err := identity.ConvertFromString(dto.Recv.A)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExprSpec(dto.Recv.ContES)
		if err != nil {
			return nil, err
		}
		return commact.RecvRec{CommChnlPH: x, ContChnlID: a, NewChnlPH: y, ContExpr: cont}, nil
	case LabExpr:
		a, err := symbol.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		contID, err := identity.ConvertFromString(dto.Lab.A)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(dto.Lab.Label)
		if err != nil {
			return nil, err
		}
		return commact.LabRec{CommChnlPH: a, ContChnlID: contID, ValExprCK: contkey.ADT(dto.Lab.ValExprCK), ValLabQN: label}, nil
	case CaseExpr:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		contID, err := identity.ConvertFromString(dto.Case.A)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]commact.ExprSpec, len(dto.Case.Branches))
		for _, branch := range dto.Case.Branches {
			cont, err := dataToExprSpec(branch.ContES)
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(branch.Label)
			if err != nil {
				return nil, err
			}
			conts[label] = cont
		}
		return commact.CaseRec{CommChnlPH: x, ContChnlID: contID, ContExprs: conts}, nil
	case FwdExpr:
		x, err := symbol.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		b, err := identity.ConvertFromString(dto.Fwd.B)
		if err != nil {
			return nil, err
		}
		return commact.FwdRec{CommChnlPH: x, ContChnlID: b}, nil
	default:
		panic(errUnexpectedExprKind(dto.K))
	}
}

func dataFromExprSpec(s commact.ExprSpec) (ExprSpec, error) {
	switch spec := s.(type) {
	case commact.CloseSpec:
		return ExprSpec{
			K:     CloseExpr,
			Close: &CloseSpec{X: symbol.ConvertToString(spec.ContChnlPH)},
		}, nil
	case commact.WaitSpec:
		dto, err := dataFromExprSpec(spec.ContExpr)
		if err != nil {
			return ExprSpec{}, err
		}
		return ExprSpec{
			K: WaitExpr,
			Wait: &WaitSpec{
				X:      symbol.ConvertToString(spec.ContChnlPH),
				ContES: dto,
			},
		}, nil
	case commact.SendSpec:
		return ExprSpec{
			K: SendExpr,
			Send: &SendSpec{
				X: symbol.ConvertToString(spec.CommChnlPH),
				Y: symbol.ConvertToString(spec.ValChnlPH),
			},
		}, nil
	case commact.RecvSpec:
		dto, err := dataFromExprSpec(spec.ContExpr)
		if err != nil {
			return ExprSpec{}, err
		}
		return ExprSpec{
			K: RecvExpr,
			Recv: &RecvSpec{
				X:      symbol.ConvertToString(spec.CommChnlPH),
				Y:      symbol.ConvertToString(spec.NewChnlPH),
				ContES: dto,
			},
		}, nil
	case commact.LabSpec:
		cont, err := dataFromExprSpec(spec.ContExpr)
		if err != nil {
			return ExprSpec{}, err
		}
		return ExprSpec{
			K: LabExpr,
			Lab: &LabSpec{
				X:      symbol.ConvertToString(spec.CommChnlPH),
				Label:  uniqsym.ConvertToString(spec.ValLabQN),
				ContES: cont,
			},
		}, nil
	case commact.CaseSpec:
		brs := []BranchSpec{}
		labels := make([]uniqsym.ADT, 0, len(spec.ContExprs))
		for label := range spec.ContExprs {
			labels = append(labels, label)
		}
		slices.SortFunc(labels, func(a, b uniqsym.ADT) int {
			return strings.Compare(uniqsym.ConvertToString(a), uniqsym.ConvertToString(b))
		})
		for _, l := range labels {
			cont := spec.ContExprs[l]
			dto, err := dataFromExprSpec(cont)
			if err != nil {
				return ExprSpec{}, err
			}
			brs = append(brs, BranchSpec{Label: uniqsym.ConvertToString(l), ContES: dto})
		}
		return ExprSpec{
			K: CaseExpr,
			Case: &CaseSpec{
				X:        symbol.ConvertToString(spec.CommChnlPH),
				Branches: brs,
			},
		}, nil
	case commact.FwdSpec:
		return ExprSpec{
			K: FwdExpr,
			Fwd: &FwdSpec{
				X: symbol.ConvertToString(spec.CommChnlPH),
				Y: symbol.ConvertToString(spec.ContChnlPH),
			},
		}, nil
	default:
		panic(commact.ErrExprTypeUnexpected(spec))
	}
}

func dataToExprSpec(dto ExprSpec) (commact.ExprSpec, error) {
	switch dto.K {
	case CloseExpr:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return commact.CloseSpec{ContChnlPH: a}, nil
	case WaitExpr:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExprSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return commact.WaitSpec{ContChnlPH: x, ContExpr: cont}, nil
	case SendExpr:
		x, err := symbol.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Send.Y)
		if err != nil {
			return nil, err
		}
		return commact.SendSpec{CommChnlPH: x, ValChnlPH: y}, nil
	case RecvExpr:
		x, err := symbol.ConvertFromString(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExprSpec(dto.Recv.ContES)
		if err != nil {
			return nil, err
		}
		return commact.RecvSpec{CommChnlPH: x, NewChnlPH: y, ContExpr: cont}, nil
	case LabExpr:
		x, err := symbol.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(dto.Lab.Label)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExprSpec(dto.Lab.ContES)
		if err != nil {
			return nil, err
		}
		return commact.LabSpec{CommChnlPH: x, ValLabQN: label, ContExpr: cont}, nil
	case CaseExpr:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]commact.ExprSpec, len(dto.Case.Branches))
		for _, b := range dto.Case.Branches {
			cont, err := dataToExprSpec(b.ContES)
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(b.Label)
			if err != nil {
				return nil, err
			}
			conts[label] = cont
		}
		return commact.CaseSpec{CommChnlPH: x, ContExprs: conts}, nil
	case FwdExpr:
		x, err := symbol.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Fwd.Y)
		if err != nil {
			return nil, err
		}
		return commact.FwdSpec{CommChnlPH: x, ContChnlPH: y}, nil
	default:
		panic(errUnexpectedExprKind(dto.K))
	}
}

func errUnexpectedExprKind(k ExprKind) error {
	return fmt.Errorf("unexpected term kind: %v", k)
}
