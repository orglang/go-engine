package adapter

import (
	"fmt"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"

	sdktermexp "github.com/orglang/go-sdk/proc/termexp"

	termexp "orglang/go-engine/proc/termexp/core"
)

func MsgFromExpSpecNilable(spec termexp.ExpSpec) *sdktermexp.ExpSpec {
	if spec == nil {
		return nil
	}
	dto := MsgFromExpSpec(spec)
	return &dto
}

func MsgFromExpSpec(s termexp.ExpSpec) sdktermexp.ExpSpec {
	switch spec := s.(type) {
	case termexp.CloseSpec:
		return sdktermexp.ExpSpec{
			K: sdktermexp.Close,
			Close: &sdktermexp.CloseSpec{
				CommChnlPH: symbol.ConvertToString(spec.ContChnlPH),
			},
		}
	case termexp.WaitSpec:
		return sdktermexp.ExpSpec{
			K: sdktermexp.Wait,
			Wait: &sdktermexp.WaitSpec{
				CommChnlPH: symbol.ConvertToString(spec.ContChnlPH),
				ContES:     MsgFromExpSpec(spec.ContExp),
			},
		}
	case termexp.SendSpec:
		return sdktermexp.ExpSpec{
			K: sdktermexp.Send,
			Send: &sdktermexp.SendSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				ValChnlPH:  symbol.ConvertToString(spec.ValChnlPH),
			},
		}
	case termexp.RecvSpec:
		return sdktermexp.ExpSpec{
			K: sdktermexp.Recv,
			Recv: &sdktermexp.RecvSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				BindChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				ContES:     MsgFromExpSpec(spec.ContExp),
			},
		}
	case termexp.LabSpec:
		return sdktermexp.ExpSpec{
			K: sdktermexp.Lab,
			Lab: &sdktermexp.LabSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				PatternQN:  uniqsym.ConvertToString(spec.ValLabQN),
			},
		}
	case termexp.CaseSpec:
		brs := []sdktermexp.BranchSpec{}
		for l, t := range spec.ContExps {
			brs = append(brs, sdktermexp.BranchSpec{PatternQN: uniqsym.ConvertToString(l), ContES: MsgFromExpSpec(t)})
		}
		return sdktermexp.ExpSpec{
			K: sdktermexp.Case,
			Case: &sdktermexp.CaseSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				ContBSes:   brs,
			},
		}
	case termexp.FwdSpec:
		return sdktermexp.ExpSpec{
			K: sdktermexp.Fwd,
			Fwd: &sdktermexp.FwdSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				ContChnlPH: symbol.ConvertToString(spec.ContChnlPH),
			},
		}
	default:
		panic(termexp.ErrExpTypeUnexpected(s))
	}
}

func MsgToExpSpecNilable(dto *sdktermexp.ExpSpec) (termexp.ExpSpec, error) {
	if dto == nil {
		return nil, nil
	}
	return MsgToExpSpec(*dto)
}

func MsgToExpSpec(dto sdktermexp.ExpSpec) (termexp.ExpSpec, error) {
	switch dto.K {
	case sdktermexp.Close:
		x, err := symbol.ConvertFromString(dto.Close.CommChnlPH)
		if err != nil {
			return nil, err
		}
		return termexp.CloseSpec{ContChnlPH: x}, nil
	case sdktermexp.Wait:
		x, err := symbol.ConvertFromString(dto.Wait.CommChnlPH)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return termexp.WaitSpec{ContChnlPH: x, ContExp: cont}, nil
	case sdktermexp.Send:
		x, err := symbol.ConvertFromString(dto.Send.CommChnlPH)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Send.ValChnlPH)
		if err != nil {
			return nil, err
		}
		return termexp.SendSpec{CommChnlPH: x, ValChnlPH: y}, nil
	case sdktermexp.Recv:
		x, err := symbol.ConvertFromString(dto.Recv.CommChnlPH)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Recv.BindChnlPH)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToExpSpec(dto.Recv.ContES)
		if err != nil {
			return nil, err
		}
		return termexp.RecvSpec{CommChnlPH: x, NewChnlPH: y, ContExp: cont}, nil
	case sdktermexp.Lab:
		x, err := symbol.ConvertFromString(dto.Lab.CommChnlPH)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(dto.Lab.PatternQN)
		if err != nil {
			return nil, err
		}
		return termexp.LabSpec{CommChnlPH: x, ValLabQN: label}, nil
	case sdktermexp.Case:
		x, err := symbol.ConvertFromString(dto.Case.CommChnlPH)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]termexp.ExpSpec, len(dto.Case.ContBSes))
		for _, b := range dto.Case.ContBSes {
			cont, err := MsgToExpSpec(b.ContES)
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(dto.Lab.PatternQN)
			if err != nil {
				return nil, err
			}
			conts[label] = cont
		}
		return termexp.CaseSpec{CommChnlPH: x, ContExps: conts}, nil
	case sdktermexp.Call:
		bindPH, err := symbol.ConvertFromString(dto.Call.BindChnlPH)
		if err != nil {
			return nil, err
		}
		procQN, err := uniqsym.ConvertFromString(dto.Call.ProcTermQN)
		if err != nil {
			return nil, err
		}
		valPHs, err := symbol.ConvertFromStrings(dto.Call.ValChnlPHs)
		if err != nil {
			return nil, err
		}
		return termexp.CallSpec{NewChnlPH: bindPH, ProcTermQN: procQN, ValChnlPHs: valPHs}, nil
	case sdktermexp.Fwd:
		x, err := symbol.ConvertFromString(dto.Fwd.CommChnlPH)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Fwd.ContChnlPH)
		if err != nil {
			return nil, err
		}
		return termexp.FwdSpec{CommChnlPH: x, ContChnlPH: y}, nil
	default:
		panic(sdktermexp.ErrUnexpectedExpKind(dto.K))
	}
}

func DataFromExpRec(r termexp.ExpRec) (termexp.ExpRecDS, error) {
	switch rec := r.(type) {
	case termexp.CloseRec:
		return termexp.ExpRecDS{
			K:     termexp.CloseExp,
			Close: &termexp.CloseRecDS{X: symbol.ConvertToString(rec.ContChnlPH)},
		}, nil
	case termexp.WaitRec:
		dto, err := dataFromExpSpec(rec.ContExp)
		if err != nil {
			return termexp.ExpRecDS{}, err
		}
		return termexp.ExpRecDS{
			K: termexp.WaitExp,
			Wait: &termexp.WaitRecDS{
				X:      symbol.ConvertToString(rec.ContChnlPH),
				ContES: dto,
			},
		}, nil
	case termexp.SendRec:
		return termexp.ExpRecDS{
			K: termexp.SendExp,
			Send: &termexp.SendRecDS{
				X: symbol.ConvertToString(rec.CommChnlPH),
				A: identity.ConvertToString(rec.ContChnlID),
				B: identity.ConvertToString(rec.ValChnlID),
			},
		}, nil
	case termexp.RecvRec:
		dto, err := dataFromExpSpec(rec.ContExp)
		if err != nil {
			return termexp.ExpRecDS{}, err
		}
		return termexp.ExpRecDS{
			K: termexp.RecvExp,
			Recv: &termexp.RecvRecDS{
				X:      symbol.ConvertToString(rec.CommChnlPH),
				Y:      symbol.ConvertToString(rec.NewChnlPH),
				ContES: dto,
			},
		}, nil
	case termexp.LabRec:
		return termexp.ExpRecDS{
			K:   termexp.LabExp,
			Lab: &termexp.LabRecDS{X: symbol.ConvertToString(rec.CommChnlPH), Label: uniqsym.ConvertToString(rec.ValLabQN)},
		}, nil
	case termexp.CaseRec:
		brs := []termexp.BranchRecDS{}
		for l, cont := range rec.ContExps {
			dto, err := dataFromExpSpec(cont)
			if err != nil {
				return termexp.ExpRecDS{}, err
			}
			brs = append(brs, termexp.BranchRecDS{Label: uniqsym.ConvertToString(l), ContES: dto})
		}
		return termexp.ExpRecDS{
			K: termexp.CaseExp,
			Case: &termexp.CaseRecDS{
				X:        symbol.ConvertToString(rec.CommChnlPH),
				Branches: brs,
			},
		}, nil
	case termexp.FwdRec:
		return termexp.ExpRecDS{
			K: termexp.FwdExp,
			Fwd: &termexp.FwdRecDS{
				X: symbol.ConvertToString(rec.CommChnlPH),
				B: identity.ConvertToString(rec.ContChnlID),
			},
		}, nil
	default:
		panic(termexp.ErrExpTypeUnexpected(rec))
	}
}

func DataToExpRec(dto termexp.ExpRecDS) (termexp.ExpRec, error) {
	switch dto.K {
	case termexp.CloseExp:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return termexp.CloseRec{ContChnlPH: a}, nil
	case termexp.WaitExp:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return termexp.WaitRec{ContChnlPH: x, ContExp: cont}, nil
	case termexp.SendExp:
		x, err := symbol.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		a, err := identity.ConvertFromString(dto.Send.A)
		if err != nil {
			return nil, err
		}
		return termexp.SendRec{CommChnlPH: x, ContChnlID: a}, nil
	case termexp.RecvExp:
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
		return termexp.RecvRec{CommChnlPH: x, NewChnlPH: y, ContExp: cont}, nil
	case termexp.LabExp:
		a, err := symbol.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(dto.Lab.Label)
		if err != nil {
			return nil, err
		}
		return termexp.LabRec{CommChnlPH: a, ValLabQN: label}, nil
	case termexp.CaseExp:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]termexp.ExpSpec, len(dto.Case.Branches))
		for _, branch := range dto.Case.Branches {
			cont, err := dataToExpSpec(branch.ContES)
			if err != nil {
				return nil, err
			}
			label, err := uniqsym.ConvertFromString(dto.Lab.Label)
			if err != nil {
				return nil, err
			}
			conts[label] = cont
		}
		return termexp.CaseRec{CommChnlPH: x, ContExps: conts}, nil
	case termexp.FwdExp:
		x, err := symbol.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		b, err := identity.ConvertFromString(dto.Fwd.B)
		if err != nil {
			return nil, err
		}
		return termexp.FwdRec{CommChnlPH: x, ContChnlID: b}, nil
	default:
		panic(errUnexpectedExpKind(dto.K))
	}
}

func dataFromExpSpec(s termexp.ExpSpec) (termexp.ExpSpecDS, error) {
	switch spec := s.(type) {
	case termexp.CloseSpec:
		return termexp.ExpSpecDS{
			K:     termexp.CloseExp,
			Close: &termexp.CloseSpecDS{X: symbol.ConvertToString(spec.ContChnlPH)},
		}, nil
	case termexp.WaitSpec:
		dto, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return termexp.ExpSpecDS{}, err
		}
		return termexp.ExpSpecDS{
			K: termexp.WaitExp,
			Wait: &termexp.WaitSpecDS{
				X:      symbol.ConvertToString(spec.ContChnlPH),
				ContES: dto,
			},
		}, nil
	case termexp.SendSpec:
		return termexp.ExpSpecDS{
			K: termexp.SendExp,
			Send: &termexp.SendSpecDS{
				X: symbol.ConvertToString(spec.CommChnlPH),
				Y: symbol.ConvertToString(spec.ValChnlPH),
			},
		}, nil
	case termexp.RecvSpec:
		dto, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return termexp.ExpSpecDS{}, err
		}
		return termexp.ExpSpecDS{
			K: termexp.RecvExp,
			Recv: &termexp.RecvSpecDS{
				X:      symbol.ConvertToString(spec.CommChnlPH),
				Y:      symbol.ConvertToString(spec.CommChnlPH),
				ContES: dto,
			},
		}, nil
	case termexp.LabSpec:
		return termexp.ExpSpecDS{
			K:   termexp.LabExp,
			Lab: &termexp.LabSpecDS{X: symbol.ConvertToString(spec.CommChnlPH), Label: uniqsym.ConvertToString(spec.ValLabQN)},
		}, nil
	case termexp.CaseSpec:
		brs := []termexp.BranchSpecDS{}
		for l, cont := range spec.ContExps {
			dto, err := dataFromExpSpec(cont)
			if err != nil {
				return termexp.ExpSpecDS{}, err
			}
			brs = append(brs, termexp.BranchSpecDS{Label: uniqsym.ConvertToString(l), ContES: dto})
		}
		return termexp.ExpSpecDS{
			K: termexp.CaseExp,
			Case: &termexp.CaseSpecDS{
				X:        symbol.ConvertToString(spec.CommChnlPH),
				Branches: brs,
			},
		}, nil
	case termexp.FwdSpec:
		return termexp.ExpSpecDS{
			K: termexp.FwdExp,
			Fwd: &termexp.FwdSpecDS{
				X: symbol.ConvertToString(spec.CommChnlPH),
				Y: symbol.ConvertToString(spec.ContChnlPH),
			},
		}, nil
	default:
		panic(termexp.ErrExpTypeUnexpected(spec))
	}
}

func dataToExpSpec(dto termexp.ExpSpecDS) (termexp.ExpSpec, error) {
	switch dto.K {
	case termexp.CloseExp:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return termexp.CloseSpec{ContChnlPH: a}, nil
	case termexp.WaitExp:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return termexp.WaitSpec{ContChnlPH: x, ContExp: cont}, nil
	case termexp.SendExp:
		x, err := symbol.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Send.Y)
		if err != nil {
			return nil, err
		}
		return termexp.SendSpec{CommChnlPH: x, ValChnlPH: y}, nil
	case termexp.RecvExp:
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
		return termexp.RecvSpec{CommChnlPH: x, NewChnlPH: y, ContExp: cont}, nil
	case termexp.LabExp:
		x, err := symbol.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(dto.Lab.Label)
		if err != nil {
			return nil, err
		}
		return termexp.LabSpec{CommChnlPH: x, ValLabQN: label}, nil
	case termexp.CaseExp:
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
		return termexp.CaseSpec{CommChnlPH: x, ContExps: conts}, nil
	case termexp.FwdExp:
		x, err := symbol.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Fwd.Y)
		if err != nil {
			return nil, err
		}
		return termexp.FwdSpec{CommChnlPH: x, ContChnlPH: y}, nil
	default:
		panic(errUnexpectedExpKind(dto.K))
	}
}

func errUnexpectedExpKind(k termexp.ExpKind) error {
	return fmt.Errorf("unexpected term kind: %v", k)
}
