package adapter

import (
	"fmt"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/proc/termexp"

	proctermexp "orglang/go-engine/proc/termexp/core"
)

func MsgFromExpSpecNilable(spec proctermexp.ExpSpec) *termexp.ExpSpec {
	if spec == nil {
		return nil
	}
	dto := MsgFromExpSpec(spec)
	return &dto
}

func MsgFromExpSpec(s proctermexp.ExpSpec) termexp.ExpSpec {
	switch spec := s.(type) {
	case proctermexp.CloseSpec:
		return termexp.ExpSpec{
			K: termexp.Close,
			Close: &termexp.CloseSpec{
				CommChnlPH: symbol.ConvertToString(spec.ContChnlPH),
			},
		}
	case proctermexp.WaitSpec:
		return termexp.ExpSpec{
			K: termexp.Wait,
			Wait: &termexp.WaitSpec{
				CommChnlPH: symbol.ConvertToString(spec.ContChnlPH),
				ContES:     MsgFromExpSpec(spec.ContExp),
			},
		}
	case proctermexp.SendSpec:
		return termexp.ExpSpec{
			K: termexp.Send,
			Send: &termexp.SendSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				ValChnlPH:  symbol.ConvertToString(spec.ValChnlPH),
			},
		}
	case proctermexp.RecvSpec:
		return termexp.ExpSpec{
			K: termexp.Recv,
			Recv: &termexp.RecvSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				BindChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				ContES:     MsgFromExpSpec(spec.ContExp),
			},
		}
	case proctermexp.LabSpec:
		return termexp.ExpSpec{
			K: termexp.Lab,
			Lab: &termexp.LabSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				PatternQN:  uniqsym.ConvertToString(spec.ValLabQN),
			},
		}
	case proctermexp.CaseSpec:
		brs := []termexp.BranchSpec{}
		for l, t := range spec.ContExps {
			brs = append(brs, termexp.BranchSpec{PatternQN: uniqsym.ConvertToString(l), ContES: MsgFromExpSpec(t)})
		}
		return termexp.ExpSpec{
			K: termexp.Case,
			Case: &termexp.CaseSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				ContBSes:   brs,
			},
		}
	case proctermexp.FwdSpec:
		return termexp.ExpSpec{
			K: termexp.Fwd,
			Fwd: &termexp.FwdSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				ContChnlPH: symbol.ConvertToString(spec.ContChnlPH),
			},
		}
	default:
		panic(proctermexp.ErrExpTypeUnexpected(s))
	}
}

func MsgToExpSpecNilable(dto *termexp.ExpSpec) (proctermexp.ExpSpec, error) {
	if dto == nil {
		return nil, nil
	}
	return MsgToExpSpec(*dto)
}

func MsgToExpSpec(dto termexp.ExpSpec) (proctermexp.ExpSpec, error) {
	switch dto.K {
	case termexp.Close:
		x, err := symbol.ConvertFromString(dto.Close.CommChnlPH)
		if err != nil {
			return nil, err
		}
		return proctermexp.CloseSpec{ContChnlPH: x}, nil
	case termexp.Wait:
		x, err := symbol.ConvertFromString(dto.Wait.CommChnlPH)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return proctermexp.WaitSpec{ContChnlPH: x, ContExp: cont}, nil
	case termexp.Send:
		x, err := symbol.ConvertFromString(dto.Send.CommChnlPH)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Send.ValChnlPH)
		if err != nil {
			return nil, err
		}
		return proctermexp.SendSpec{CommChnlPH: x, ValChnlPH: y}, nil
	case termexp.Recv:
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
		return proctermexp.RecvSpec{CommChnlPH: x, NewChnlPH: y, ContExp: cont}, nil
	case termexp.Lab:
		x, err := symbol.ConvertFromString(dto.Lab.CommChnlPH)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(dto.Lab.PatternQN)
		if err != nil {
			return nil, err
		}
		return proctermexp.LabSpec{CommChnlPH: x, ValLabQN: label}, nil
	case termexp.Case:
		x, err := symbol.ConvertFromString(dto.Case.CommChnlPH)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]proctermexp.ExpSpec, len(dto.Case.ContBSes))
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
		return proctermexp.CaseSpec{CommChnlPH: x, ContExps: conts}, nil
	case termexp.Call:
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
		return proctermexp.CallSpec{NewChnlPH: bindPH, ProcTermQN: procQN, ValChnlPHs: valPHs}, nil
	case termexp.Fwd:
		x, err := symbol.ConvertFromString(dto.Fwd.CommChnlPH)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Fwd.ContChnlPH)
		if err != nil {
			return nil, err
		}
		return proctermexp.FwdSpec{CommChnlPH: x, ContChnlPH: y}, nil
	default:
		panic(termexp.ErrUnexpectedExpKind(dto.K))
	}
}

func DataFromExpRec(r proctermexp.ExpRec) (proctermexp.ExpRecDS, error) {
	switch rec := r.(type) {
	case proctermexp.CloseRec:
		return proctermexp.ExpRecDS{
			K:     proctermexp.CloseExp,
			Close: &proctermexp.CloseRecDS{X: symbol.ConvertToString(rec.ContChnlPH)},
		}, nil
	case proctermexp.WaitRec:
		dto, err := dataFromExpSpec(rec.ContExp)
		if err != nil {
			return proctermexp.ExpRecDS{}, err
		}
		return proctermexp.ExpRecDS{
			K: proctermexp.WaitExp,
			Wait: &proctermexp.WaitRecDS{
				X:      symbol.ConvertToString(rec.ContChnlPH),
				ContES: dto,
			},
		}, nil
	case proctermexp.SendRec:
		return proctermexp.ExpRecDS{
			K: proctermexp.SendExp,
			Send: &proctermexp.SendRecDS{
				X: symbol.ConvertToString(rec.CommChnlPH),
				A: identity.ConvertToString(rec.ContChnlID),
				B: identity.ConvertToString(rec.ValChnlID),
			},
		}, nil
	case proctermexp.RecvRec:
		dto, err := dataFromExpSpec(rec.ContExp)
		if err != nil {
			return proctermexp.ExpRecDS{}, err
		}
		return proctermexp.ExpRecDS{
			K: proctermexp.RecvExp,
			Recv: &proctermexp.RecvRecDS{
				X:      symbol.ConvertToString(rec.CommChnlPH),
				Y:      symbol.ConvertToString(rec.NewChnlPH),
				ContES: dto,
			},
		}, nil
	case proctermexp.LabRec:
		return proctermexp.ExpRecDS{
			K:   proctermexp.LabExp,
			Lab: &proctermexp.LabRecDS{X: symbol.ConvertToString(rec.CommChnlPH), Label: uniqsym.ConvertToString(rec.ValLabQN)},
		}, nil
	case proctermexp.CaseRec:
		brs := []proctermexp.BranchRecDS{}
		for l, cont := range rec.ContExps {
			dto, err := dataFromExpSpec(cont)
			if err != nil {
				return proctermexp.ExpRecDS{}, err
			}
			brs = append(brs, proctermexp.BranchRecDS{Label: uniqsym.ConvertToString(l), ContES: dto})
		}
		return proctermexp.ExpRecDS{
			K: proctermexp.CaseExp,
			Case: &proctermexp.CaseRecDS{
				X:        symbol.ConvertToString(rec.CommChnlPH),
				Branches: brs,
			},
		}, nil
	case proctermexp.FwdRec:
		return proctermexp.ExpRecDS{
			K: proctermexp.FwdExp,
			Fwd: &proctermexp.FwdRecDS{
				X: symbol.ConvertToString(rec.CommChnlPH),
				B: identity.ConvertToString(rec.ContChnlID),
			},
		}, nil
	default:
		panic(proctermexp.ErrExpTypeUnexpected(rec))
	}
}

func DataToExpRec(dto proctermexp.ExpRecDS) (proctermexp.ExpRec, error) {
	switch dto.K {
	case proctermexp.CloseExp:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return proctermexp.CloseRec{ContChnlPH: a}, nil
	case proctermexp.WaitExp:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return proctermexp.WaitRec{ContChnlPH: x, ContExp: cont}, nil
	case proctermexp.SendExp:
		x, err := symbol.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		a, err := identity.ConvertFromString(dto.Send.A)
		if err != nil {
			return nil, err
		}
		return proctermexp.SendRec{CommChnlPH: x, ContChnlID: a}, nil
	case proctermexp.RecvExp:
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
		return proctermexp.RecvRec{CommChnlPH: x, NewChnlPH: y, ContExp: cont}, nil
	case proctermexp.LabExp:
		a, err := symbol.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(dto.Lab.Label)
		if err != nil {
			return nil, err
		}
		return proctermexp.LabRec{CommChnlPH: a, ValLabQN: label}, nil
	case proctermexp.CaseExp:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]proctermexp.ExpSpec, len(dto.Case.Branches))
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
		return proctermexp.CaseRec{CommChnlPH: x, ContExps: conts}, nil
	case proctermexp.FwdExp:
		x, err := symbol.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		b, err := identity.ConvertFromString(dto.Fwd.B)
		if err != nil {
			return nil, err
		}
		return proctermexp.FwdRec{CommChnlPH: x, ContChnlID: b}, nil
	default:
		panic(errUnexpectedExpKind(dto.K))
	}
}

func dataFromExpSpec(s proctermexp.ExpSpec) (proctermexp.ExpSpecDS, error) {
	switch spec := s.(type) {
	case proctermexp.CloseSpec:
		return proctermexp.ExpSpecDS{
			K:     proctermexp.CloseExp,
			Close: &proctermexp.CloseSpecDS{X: symbol.ConvertToString(spec.ContChnlPH)},
		}, nil
	case proctermexp.WaitSpec:
		dto, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return proctermexp.ExpSpecDS{}, err
		}
		return proctermexp.ExpSpecDS{
			K: proctermexp.WaitExp,
			Wait: &proctermexp.WaitSpecDS{
				X:      symbol.ConvertToString(spec.ContChnlPH),
				ContES: dto,
			},
		}, nil
	case proctermexp.SendSpec:
		return proctermexp.ExpSpecDS{
			K: proctermexp.SendExp,
			Send: &proctermexp.SendSpecDS{
				X: symbol.ConvertToString(spec.CommChnlPH),
				Y: symbol.ConvertToString(spec.ValChnlPH),
			},
		}, nil
	case proctermexp.RecvSpec:
		dto, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return proctermexp.ExpSpecDS{}, err
		}
		return proctermexp.ExpSpecDS{
			K: proctermexp.RecvExp,
			Recv: &proctermexp.RecvSpecDS{
				X:      symbol.ConvertToString(spec.CommChnlPH),
				Y:      symbol.ConvertToString(spec.CommChnlPH),
				ContES: dto,
			},
		}, nil
	case proctermexp.LabSpec:
		return proctermexp.ExpSpecDS{
			K:   proctermexp.LabExp,
			Lab: &proctermexp.LabSpecDS{X: symbol.ConvertToString(spec.CommChnlPH), Label: uniqsym.ConvertToString(spec.ValLabQN)},
		}, nil
	case proctermexp.CaseSpec:
		brs := []proctermexp.BranchSpecDS{}
		for l, cont := range spec.ContExps {
			dto, err := dataFromExpSpec(cont)
			if err != nil {
				return proctermexp.ExpSpecDS{}, err
			}
			brs = append(brs, proctermexp.BranchSpecDS{Label: uniqsym.ConvertToString(l), ContES: dto})
		}
		return proctermexp.ExpSpecDS{
			K: proctermexp.CaseExp,
			Case: &proctermexp.CaseSpecDS{
				X:        symbol.ConvertToString(spec.CommChnlPH),
				Branches: brs,
			},
		}, nil
	case proctermexp.FwdSpec:
		return proctermexp.ExpSpecDS{
			K: proctermexp.FwdExp,
			Fwd: &proctermexp.FwdSpecDS{
				X: symbol.ConvertToString(spec.CommChnlPH),
				Y: symbol.ConvertToString(spec.ContChnlPH),
			},
		}, nil
	default:
		panic(proctermexp.ErrExpTypeUnexpected(spec))
	}
}

func dataToExpSpec(dto proctermexp.ExpSpecDS) (proctermexp.ExpSpec, error) {
	switch dto.K {
	case proctermexp.CloseExp:
		a, err := symbol.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return proctermexp.CloseSpec{ContChnlPH: a}, nil
	case proctermexp.WaitExp:
		x, err := symbol.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToExpSpec(dto.Wait.ContES)
		if err != nil {
			return nil, err
		}
		return proctermexp.WaitSpec{ContChnlPH: x, ContExp: cont}, nil
	case proctermexp.SendExp:
		x, err := symbol.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Send.Y)
		if err != nil {
			return nil, err
		}
		return proctermexp.SendSpec{CommChnlPH: x, ValChnlPH: y}, nil
	case proctermexp.RecvExp:
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
		return proctermexp.RecvSpec{CommChnlPH: x, NewChnlPH: y, ContExp: cont}, nil
	case proctermexp.LabExp:
		x, err := symbol.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(dto.Lab.Label)
		if err != nil {
			return nil, err
		}
		return proctermexp.LabSpec{CommChnlPH: x, ValLabQN: label}, nil
	case proctermexp.CaseExp:
		x, err := symbol.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[uniqsym.ADT]proctermexp.ExpSpec, len(dto.Case.Branches))
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
		return proctermexp.CaseSpec{CommChnlPH: x, ContExps: conts}, nil
	case proctermexp.FwdExp:
		x, err := symbol.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		y, err := symbol.ConvertFromString(dto.Fwd.Y)
		if err != nil {
			return nil, err
		}
		return proctermexp.FwdSpec{CommChnlPH: x, ContChnlPH: y}, nil
	default:
		panic(errUnexpectedExpKind(dto.K))
	}
}

func errUnexpectedExpKind(k proctermexp.ExpKind) error {
	return fmt.Errorf("unexpected term kind: %v", k)
}
