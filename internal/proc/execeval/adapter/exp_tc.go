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

func DataFromExpRec(r commact.ExpRec) (ExpRec, error) {
	switch rec := r.(type) {
	case commact.CloseRec:
		return ExpRec{
			K:     CloseExp,
			Close: &CloseRec{X: symbol.ConvertToString(rec.ContChnlPH)},
		}, nil
	case commact.WaitRec:
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
	case commact.SendRec:
		commID := identity.ConvertToString(rec.CommID)
		return ExpRec{
			K: SendExp,
			Send: &SendRec{
				X:        symbol.ConvertToString(rec.CommChnlPH),
				A:        identity.ConvertToString(rec.ContChnlID),
				B:        identity.ConvertToString(rec.ValChnlID),
				CommID:   commID,
				ValExpCK: int64(rec.ValExpCK),
			},
		}, nil
	case commact.RecvRec:
		dto, err := dataFromExpSpec(rec.ContExp)
		if err != nil {
			return ExpRec{}, err
		}
		return ExpRec{
			K: RecvExp,
			Recv: &RecvRec{
				X:      symbol.ConvertToString(rec.CommChnlPH),
				A:      identity.ConvertToString(rec.ContChnlID),
				Y:      symbol.ConvertToString(rec.NewChnlPH),
				ContES: dto,
			},
		}, nil
	case commact.LabRec:
		return ExpRec{
			K: LabExp,
			Lab: &LabRec{
				X:        symbol.ConvertToString(rec.CommChnlPH),
				A:        identity.ConvertToString(rec.ContChnlID),
				Label:    uniqsym.ConvertToString(rec.ValLabQN),
				ValExpCK: int64(rec.ValExpCK),
			},
		}, nil
	case commact.CaseRec:
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
				X:        symbol.ConvertToString(rec.CommChnlPH),
				A:        identity.ConvertToString(rec.ContChnlID),
				Branches: brs,
			},
		}, nil
	case commact.FwdRec:
		return ExpRec{
			K: FwdExp,
			Fwd: &FwdRec{
				X: symbol.ConvertToString(rec.CommChnlPH),
				B: identity.ConvertToString(rec.ContChnlID),
			},
		}, nil
	default:
		panic(commact.ErrExpTypeUnexpected(rec))
	}
}

func DataToExpRec(dto ExpRec) (commact.ExpRec, error) {
	switch dto.K {
	case CloseExp:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return commact.CloseRec{ContChnlPH: a}, nil
	case WaitExp:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return commact.WaitRec{ContChnlPH: x, ContExp: cont}, nil
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
		commID, err := identity.ConvertFromString(dto.Send.CommID)
		if err != nil {
			return nil, err
		}
		return commact.SendRec{
			CommChnlPH: x,
			ContChnlID: a,
			ValChnlID:  valID,
			ValExpCK:   contkey.ADT(dto.Send.ValExpCK),
			CommID:     commID,
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
		return commact.RecvRec{CommChnlPH: x, ContChnlID: a, NewChnlPH: y, ContExp: cont}, nil
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
		return commact.LabRec{CommChnlPH: a, ContChnlID: contID, ValExpCK: contkey.ADT(dto.Lab.ValExpCK), ValLabQN: label}, nil
	case CaseExp:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		contID, err := identity.ConvertFromString(dto.Case.A)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]commact.ExpSpec, len(dto.Case.Branches))
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
		return commact.CaseRec{CommChnlPH: x, ContChnlID: contID, ContExps: conts}, nil
	case FwdExp:
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
		panic(errUnexpectedExpKind(dto.K))
	}
}

func dataFromExpSpec(s commact.ExpSpec) (ExpSpec, error) {
	switch spec := s.(type) {
	case commact.CloseSpec:
		return ExpSpec{
			K:     CloseExp,
			Close: &CloseSpec{X: symbol.ConvertToString(spec.ContChnlPH)},
		}, nil
	case commact.WaitSpec:
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
	case commact.SendSpec:
		return ExpSpec{
			K: SendExp,
			Send: &SendSpec{
				X: symbol.ConvertToString(spec.CommChnlPH),
				Y: symbol.ConvertToString(spec.ValChnlPH),
			},
		}, nil
	case commact.RecvSpec:
		dto, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return ExpSpec{}, err
		}
		return ExpSpec{
			K: RecvExp,
			Recv: &RecvSpec{
				X:      symbol.ConvertToString(spec.CommChnlPH),
				Y:      symbol.ConvertToString(spec.NewChnlPH),
				ContES: dto,
			},
		}, nil
	case commact.LabSpec:
		cont, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return ExpSpec{}, err
		}
		return ExpSpec{
			K: LabExp,
			Lab: &LabSpec{
				X:      symbol.ConvertToString(spec.CommChnlPH),
				Label:  uniqsym.ConvertToString(spec.ValLabQN),
				ContES: cont,
			},
		}, nil
	case commact.CaseSpec:
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
				X:        symbol.ConvertToString(spec.CommChnlPH),
				Branches: brs,
			},
		}, nil
	case commact.FwdSpec:
		return ExpSpec{
			K: FwdExp,
			Fwd: &FwdSpec{
				X: symbol.ConvertToString(spec.CommChnlPH),
				Y: symbol.ConvertToString(spec.ContChnlPH),
			},
		}, nil
	default:
		panic(commact.ErrExpTypeUnexpected(spec))
	}
}

func dataToExpSpec(dto ExpSpec) (commact.ExpSpec, error) {
	switch dto.K {
	case CloseExp:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return commact.CloseSpec{ContChnlPH: a}, nil
	case WaitExp:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return commact.WaitSpec{ContChnlPH: x, ContExp: cont}, nil
	case SendExp:
		x, err := symbol.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Send.Y)
		if err != nil {
			return nil, err
		}
		return commact.SendSpec{CommChnlPH: x, ValChnlPH: y}, nil
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
		return commact.RecvSpec{CommChnlPH: x, NewChnlPH: y, ContExp: cont}, nil
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
		return commact.LabSpec{CommChnlPH: x, ValLabQN: label, ContExp: cont}, nil
	case CaseExp:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]commact.ExpSpec, len(dto.Case.Branches))
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
		return commact.CaseSpec{CommChnlPH: x, ContExps: conts}, nil
	case FwdExp:
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
		panic(errUnexpectedExpKind(dto.K))
	}
}

func errUnexpectedExpKind(k ExpKind) error {
	return fmt.Errorf("unexpected term kind: %v", k)
}
