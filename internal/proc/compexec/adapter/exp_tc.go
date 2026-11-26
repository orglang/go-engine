package adapter

import (
	"fmt"
	"slices"
	"strings"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	compexp "orglang/go-engine/internal/proc/compexec/core"
)

func DataFromExpRec(r compexp.ExpRec) (ExpRec, error) {
	switch rec := r.(type) {
	case compexp.CloseRec:
		return ExpRec{
			K:     CloseExp,
			Close: &CloseRec{X: symbol.ConvertToString(rec.ContChnlPH)},
		}, nil
	case compexp.WaitRec:
		dto, err := dataFromExpSpec(rec.ContExp)
		if err != nil {
			return ExpRec{}, err
		}
		return ExpRec{
			K: WaitExp,
			Wait: &WaitRec{
				X:      symbol.ConvertToString(rec.ContChnlPH),
				ContES: dto,
			},
		}, nil
	case compexp.SendRec:
		connID := identity.ConvertToString(rec.ConnID)
		return ExpRec{
			K: SendExp,
			Send: &SendRec{
				X:        symbol.ConvertToString(rec.ConnChnlPH),
				A:        identity.ConvertToString(rec.ContChnlID),
				B:        identity.ConvertToString(rec.ValChnlID),
				ConnID:   connID,
				ValExpCK: int64(rec.ValExpCK),
			},
		}, nil
	case compexp.RecvRec:
		dto, err := dataFromExpSpec(rec.ContExp)
		if err != nil {
			return ExpRec{}, err
		}
		return ExpRec{
			K: RecvExp,
			Recv: &RecvRec{
				X:      symbol.ConvertToString(rec.ConnChnlPH),
				A:      identity.ConvertToString(rec.ContChnlID),
				Y:      symbol.ConvertToString(rec.NewChnlPH),
				ContES: dto,
			},
		}, nil
	case compexp.LabRec:
		return ExpRec{
			K: LabExp,
			Lab: &LabRec{
				X:        symbol.ConvertToString(rec.ConnChnlPH),
				A:        identity.ConvertToString(rec.ContChnlID),
				Label:    uniqsym.ConvertToString(rec.ValLabQN),
				ValExpCK: int64(rec.ValExpCK),
			},
		}, nil
	case compexp.CaseRec:
		brs := []BranchRec{}
		for l, cont := range rec.ContExps {
			dto, err := dataFromExpSpec(cont)
			if err != nil {
				return ExpRec{}, err
			}
			brs = append(brs, BranchRec{Label: uniqsym.ConvertToString(l), ContES: dto})
		}
		return ExpRec{
			K: CaseExp,
			Case: &CaseRec{
				X:        symbol.ConvertToString(rec.ConnChnlPH),
				A:        identity.ConvertToString(rec.ContChnlID),
				Branches: brs,
			},
		}, nil
	case compexp.FwdRec:
		return ExpRec{
			K: FwdExp,
			Fwd: &FwdRec{
				X: symbol.ConvertToString(rec.ConnChnlPH),
				B: identity.ConvertToString(rec.ContChnlID),
			},
		}, nil
	default:
		panic(compexp.ErrExpTypeUnexpected(rec))
	}
}

func DataToExpRec(dto ExpRec) (compexp.ExpRec, error) {
	switch dto.K {
	case CloseExp:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return compexp.CloseRec{ContChnlPH: a}, nil
	case WaitExp:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return compexp.WaitRec{ContChnlPH: x, ContExp: cont}, nil
	case SendExp:
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
		connID, err := identity.ConvertFromString(dto.Send.ConnID)
		if err != nil {
			return nil, err
		}
		return compexp.SendRec{
			ConnChnlPH: x,
			ContChnlID: a,
			ValChnlID:  valID,
			ValExpCK:   contkey.ADT(dto.Send.ValExpCK),
			ConnID:     connID,
		}, nil
	case RecvExp:
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
		cont, err := dataToExpSpec(dto.Recv.ContES)
		if err != nil {
			return nil, err
		}
		return compexp.RecvRec{ConnChnlPH: x, ContChnlID: a, NewChnlPH: y, ContExp: cont}, nil
	case LabExp:
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
		return compexp.LabRec{ConnChnlPH: a, ContChnlID: contID, ValExpCK: contkey.ADT(dto.Lab.ValExpCK), ValLabQN: label}, nil
	case CaseExp:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		contID, err := identity.ConvertFromString(dto.Case.A)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]compexp.ExpSpec, len(dto.Case.Branches))
		for _, branch := range dto.Case.Branches {
			cont, err := dataToExpSpec(branch.ContES)
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(branch.Label)
			if err != nil {
				return nil, err
			}
			conts[label] = cont
		}
		return compexp.CaseRec{ConnChnlPH: x, ContChnlID: contID, ContExps: conts}, nil
	case FwdExp:
		x, err := symbol.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		b, err := identity.ConvertFromString(dto.Fwd.B)
		if err != nil {
			return nil, err
		}
		return compexp.FwdRec{ConnChnlPH: x, ContChnlID: b}, nil
	default:
		panic(errUnexpectedExpKind(dto.K))
	}
}

func dataFromExpSpec(s compexp.ExpSpec) (ExpSpec, error) {
	switch spec := s.(type) {
	case compexp.CloseSpec:
		return ExpSpec{
			K:     CloseExp,
			Close: &CloseSpec{X: symbol.ConvertToString(spec.ContChnlPH)},
		}, nil
	case compexp.WaitSpec:
		dto, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return ExpSpec{}, err
		}
		return ExpSpec{
			K: WaitExp,
			Wait: &WaitSpec{
				X:      symbol.ConvertToString(spec.ContChnlPH),
				ContES: dto,
			},
		}, nil
	case compexp.SendSpec:
		return ExpSpec{
			K: SendExp,
			Send: &SendSpec{
				X: symbol.ConvertToString(spec.ConnChnlPH),
				Y: symbol.ConvertToString(spec.ValChnlPH),
			},
		}, nil
	case compexp.RecvSpec:
		dto, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return ExpSpec{}, err
		}
		return ExpSpec{
			K: RecvExp,
			Recv: &RecvSpec{
				X:      symbol.ConvertToString(spec.ConnChnlPH),
				Y:      symbol.ConvertToString(spec.NewChnlPH),
				ContES: dto,
			},
		}, nil
	case compexp.LabSpec:
		cont, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return ExpSpec{}, err
		}
		return ExpSpec{
			K: LabExp,
			Lab: &LabSpec{
				X:      symbol.ConvertToString(spec.ConnChnlPH),
				Label:  uniqsym.ConvertToString(spec.ValLabQN),
				ContES: cont,
			},
		}, nil
	case compexp.CaseSpec:
		brs := []BranchSpec{}
		labels := make([]uniqsym.ADT, 0, len(spec.ContExps))
		for label := range spec.ContExps {
			labels = append(labels, label)
		}
		slices.SortFunc(labels, func(a, b uniqsym.ADT) int {
			return strings.Compare(uniqsym.ConvertToString(a), uniqsym.ConvertToString(b))
		})
		for _, l := range labels {
			cont := spec.ContExps[l]
			dto, err := dataFromExpSpec(cont)
			if err != nil {
				return ExpSpec{}, err
			}
			brs = append(brs, BranchSpec{Label: uniqsym.ConvertToString(l), ContES: dto})
		}
		return ExpSpec{
			K: CaseExp,
			Case: &CaseSpec{
				X:        symbol.ConvertToString(spec.ConnChnlPH),
				Branches: brs,
			},
		}, nil
	case compexp.FwdSpec:
		return ExpSpec{
			K: FwdExp,
			Fwd: &FwdSpec{
				X: symbol.ConvertToString(spec.ConnChnlPH),
				Y: symbol.ConvertToString(spec.ContChnlPH),
			},
		}, nil
	default:
		panic(compexp.ErrExpTypeUnexpected(spec))
	}
}

func dataToExpSpec(dto ExpSpec) (compexp.ExpSpec, error) {
	switch dto.K {
	case CloseExp:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return compexp.CloseSpec{ContChnlPH: a}, nil
	case WaitExp:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return compexp.WaitSpec{ContChnlPH: x, ContExp: cont}, nil
	case SendExp:
		x, err := symbol.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Send.Y)
		if err != nil {
			return nil, err
		}
		return compexp.SendSpec{ConnChnlPH: x, ValChnlPH: y}, nil
	case RecvExp:
		x, err := symbol.ConvertFromString(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Recv.ContES)
		if err != nil {
			return nil, err
		}
		return compexp.RecvSpec{ConnChnlPH: x, NewChnlPH: y, ContExp: cont}, nil
	case LabExp:
		x, err := symbol.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(dto.Lab.Label)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Lab.ContES)
		if err != nil {
			return nil, err
		}
		return compexp.LabSpec{ConnChnlPH: x, ValLabQN: label, ContExp: cont}, nil
	case CaseExp:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]compexp.ExpSpec, len(dto.Case.Branches))
		for _, b := range dto.Case.Branches {
			cont, err := dataToExpSpec(b.ContES)
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(b.Label)
			if err != nil {
				return nil, err
			}
			conts[label] = cont
		}
		return compexp.CaseSpec{ConnChnlPH: x, ContExps: conts}, nil
	case FwdExp:
		x, err := symbol.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Fwd.Y)
		if err != nil {
			return nil, err
		}
		return compexp.FwdSpec{ConnChnlPH: x, ContChnlPH: y}, nil
	default:
		panic(errUnexpectedExpKind(dto.K))
	}
}

func errUnexpectedExpKind(k ExpKind) error {
	return fmt.Errorf("unexpected term kind: %v", k)
}
