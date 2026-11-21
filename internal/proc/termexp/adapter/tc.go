package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	termexp "orglang/go-engine/internal/proc/termexp/core"
)

func DataFromExpRec(r termexp.ExpRec) (ExpRec, error) {
	switch rec := r.(type) {
	case termexp.CloseRec:
		return ExpRec{
			K:     CloseExp,
			Close: &CloseRec{X: symbol.ConvertToString(rec.ContChnlPH)},
		}, nil
	case termexp.WaitRec:
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
	case termexp.SendRec:
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
	case termexp.RecvRec:
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
	case termexp.LabRec:
		return ExpRec{
			K: LabExp,
			Lab: &LabRec{
				X:        symbol.ConvertToString(rec.ConnChnlPH),
				A:        identity.ConvertToString(rec.ContChnlID),
				Label:    uniqsym.ConvertToString(rec.ValLabQN),
				ValExpCK: int64(rec.ValExpCK),
			},
		}, nil
	case termexp.CaseRec:
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
	case termexp.FwdRec:
		return ExpRec{
			K: FwdExp,
			Fwd: &FwdRec{
				X: symbol.ConvertToString(rec.ConnChnlPH),
				B: identity.ConvertToString(rec.ContChnlID),
			},
		}, nil
	default:
		panic(termexp.ErrExpTypeUnexpected(rec))
	}
}

func DataToExpRec(dto ExpRec) (termexp.ExpRec, error) {
	switch dto.K {
	case CloseExp:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return termexp.CloseRec{ContChnlPH: a}, nil
	case WaitExp:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return termexp.WaitRec{ContChnlPH: x, ContExp: cont}, nil
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
		return termexp.SendRec{
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
		return termexp.RecvRec{ConnChnlPH: x, ContChnlID: a, NewChnlPH: y, ContExp: cont}, nil
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
		return termexp.LabRec{ConnChnlPH: a, ContChnlID: contID, ValExpCK: contkey.ADT(dto.Lab.ValExpCK), ValLabQN: label}, nil
	case CaseExp:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		contID, err := identity.ConvertFromString(dto.Case.A)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]termexp.ExpSpec, len(dto.Case.Branches))
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
		return termexp.CaseRec{ConnChnlPH: x, ContChnlID: contID, ContExps: conts}, nil
	case FwdExp:
		x, err := symbol.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		b, err := identity.ConvertFromString(dto.Fwd.B)
		if err != nil {
			return nil, err
		}
		return termexp.FwdRec{ConnChnlPH: x, ContChnlID: b}, nil
	default:
		panic(errUnexpectedExpKind(dto.K))
	}
}

func dataFromExpSpec(s termexp.ExpSpec) (ExpSpec, error) {
	switch spec := s.(type) {
	case termexp.CloseSpec:
		return ExpSpec{
			K:     CloseExp,
			Close: &CloseSpec{X: symbol.ConvertToString(spec.ContChnlPH)},
		}, nil
	case termexp.WaitSpec:
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
	case termexp.SendSpec:
		return ExpSpec{
			K: SendExp,
			Send: &SendSpec{
				X: symbol.ConvertToString(spec.ConnChnlPH),
				Y: symbol.ConvertToString(spec.ValChnlPH),
			},
		}, nil
	case termexp.RecvSpec:
		dto, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return ExpSpec{}, err
		}
		return ExpSpec{
			K: RecvExp,
			Recv: &RecvSpec{
				X:      symbol.ConvertToString(spec.ConnChnlPH),
				Y:      symbol.ConvertToString(spec.ConnChnlPH),
				ContES: dto,
			},
		}, nil
	case termexp.LabSpec:
		return ExpSpec{
			K:   LabExp,
			Lab: &LabSpec{X: symbol.ConvertToString(spec.ConnChnlPH), Label: uniqsym.ConvertToString(spec.ValLabQN)},
		}, nil
	case termexp.CaseSpec:
		brs := []BranchSpec{}
		for l, cont := range spec.ContExps {
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
	case termexp.FwdSpec:
		return ExpSpec{
			K: FwdExp,
			Fwd: &FwdSpec{
				X: symbol.ConvertToString(spec.ConnChnlPH),
				Y: symbol.ConvertToString(spec.ContChnlPH),
			},
		}, nil
	default:
		panic(termexp.ErrExpTypeUnexpected(spec))
	}
}

func dataToExpSpec(dto ExpSpec) (termexp.ExpSpec, error) {
	switch dto.K {
	case CloseExp:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return termexp.CloseSpec{ContChnlPH: a}, nil
	case WaitExp:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return termexp.WaitSpec{ContChnlPH: x, ContExp: cont}, nil
	case SendExp:
		x, err := symbol.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Send.Y)
		if err != nil {
			return nil, err
		}
		return termexp.SendSpec{ConnChnlPH: x, ValChnlPH: y}, nil
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
		return termexp.RecvSpec{ConnChnlPH: x, NewChnlPH: y, ContExp: cont}, nil
	case LabExp:
		x, err := symbol.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(dto.Lab.Label)
		if err != nil {
			return nil, err
		}
		return termexp.LabSpec{ConnChnlPH: x, ValLabQN: label}, nil
	case CaseExp:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]termexp.ExpSpec, len(dto.Case.Branches))
		for _, b := range dto.Case.Branches {
			cont, err := dataToExpSpec(b.ContES)
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(dto.Lab.Label)
			if err != nil {
				return nil, err
			}
			conts[label] = cont
		}
		return termexp.CaseSpec{ConnChnlPH: x, ContExps: conts}, nil
	case FwdExp:
		x, err := symbol.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Fwd.Y)
		if err != nil {
			return nil, err
		}
		return termexp.FwdSpec{ConnChnlPH: x, ContChnlPH: y}, nil
	default:
		panic(errUnexpectedExpKind(dto.K))
	}
}

func errUnexpectedExpKind(k ExpKind) error {
	return fmt.Errorf("unexpected term kind: %v", k)
}
