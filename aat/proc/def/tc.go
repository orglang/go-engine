package def

import (
	"fmt"

	"orglang/orglang/avt/id"
	"orglang/orglang/avt/sym"
)

func MsgFromTermSpecNilable(t TermSpec) *TermSpecME {
	if t == nil {
		return nil
	}
	term := MsgFromTermSpec(t)
	return &term
}

func MsgFromTermSpec(t TermSpec) TermSpecME {
	switch term := t.(type) {
	case CloseSpec:
		return TermSpecME{
			K: Close,
			Close: &CloseSpecME{
				X: sym.ConvertToString(term.CommPH),
			},
		}
	case WaitSpec:
		return TermSpecME{
			K: Wait,
			Wait: &WaitSpecME{
				X:    sym.ConvertToString(term.CommPH),
				Cont: MsgFromTermSpec(term.ContTS),
			},
		}
	case SendSpec:
		return TermSpecME{
			K: Send,
			Send: &SendSpecME{
				X: sym.ConvertToString(term.CommPH),
				Y: sym.ConvertToString(term.ValPH),
			},
		}
	case RecvSpec:
		return TermSpecME{
			K: Recv,
			Recv: &RecvSpecME{
				X:    sym.ConvertToString(term.CommPH),
				Y:    sym.ConvertToString(term.BindPH),
				Cont: MsgFromTermSpec(term.ContTS),
			},
		}
	case LabSpec:
		return TermSpecME{
			K: Lab,
			Lab: &LabSpecME{
				X:     sym.ConvertToString(term.CommPH),
				Label: string(term.Label),
			},
		}
	case CaseSpec:
		brs := []BranchSpecME{}
		for l, t := range term.Conts {
			brs = append(brs, BranchSpecME{Label: string(l), Cont: MsgFromTermSpec(t)})
		}
		return TermSpecME{
			K: Case,
			Case: &CaseSpecME{
				X:   sym.ConvertToString(term.CommPH),
				Brs: brs,
			},
		}
	case SpawnSpecOld:
		return TermSpecME{
			K: Spawn,
			Spawn: &SpawnSpecME{
				X:     sym.ConvertToString(term.X),
				SigID: id.ConvertToString(term.SigID),
				Ys:    sym.ConvertToStrings(term.Ys),
				Cont:  MsgFromTermSpecNilable(term.Cont),
			},
		}
	case CallSpecOld:
		return TermSpecME{
			K: Call,
			Call: &CallSpecME{
				X:     sym.ConvertToString(term.X),
				SigPH: sym.ConvertToString(term.SigPH),
				Ys:    sym.ConvertToStrings(term.Ys),
			},
		}
	case FwdSpec:
		return TermSpecME{
			K: Fwd,
			Fwd: &FwdSpecME{
				X: sym.ConvertToString(term.X),
				Y: sym.ConvertToString(term.Y),
			},
		}
	default:
		panic(ErrTermTypeUnexpected(t))
	}
}

func MsgToTermSpecNilable(dto *TermSpecME) (TermSpec, error) {
	if dto == nil {
		return nil, nil
	}
	return MsgToTermSpec(*dto)
}

func MsgToTermSpec(dto TermSpecME) (TermSpec, error) {
	switch dto.K {
	case Close:
		x, err := sym.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return CloseSpec{CommPH: x}, nil
	case Wait:
		x, err := sym.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTermSpec(dto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitSpec{CommPH: x, ContTS: cont}, nil
	case Send:
		x, err := sym.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		y, err := sym.ConvertFromString(dto.Send.Y)
		if err != nil {
			return nil, err
		}
		return SendSpec{CommPH: x, ValPH: y}, nil
	case Recv:
		x, err := sym.ConvertFromString(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := sym.ConvertFromString(dto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTermSpec(dto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return RecvSpec{CommPH: x, BindPH: y, ContTS: cont}, nil
	case Lab:
		x, err := sym.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		return LabSpec{CommPH: x, Label: sym.ADT(dto.Lab.Label)}, nil
	case Case:
		x, err := sym.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[sym.ADT]TermSpec, len(dto.Case.Brs))
		for _, b := range dto.Case.Brs {
			cont, err := MsgToTermSpec(b.Cont)
			if err != nil {
				return nil, err
			}
			conts[sym.ADT(b.Label)] = cont
		}
		return CaseSpec{CommPH: x, Conts: conts}, nil
	case Spawn:
		x, err := sym.ConvertFromString(dto.Spawn.X)
		if err != nil {
			return nil, err
		}
		sigID, err := id.ConvertFromString(dto.Spawn.SigID)
		if err != nil {
			return nil, err
		}
		ys, err := sym.ConvertFromStrings(dto.Spawn.Ys)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTermSpecNilable(dto.Spawn.Cont)
		if err != nil {
			return nil, err
		}
		return SpawnSpecOld{X: x, SigID: sigID, Ys: ys, Cont: cont}, nil
	case Call:
		x, err := sym.ConvertFromString(dto.Call.X)
		if err != nil {
			return nil, err
		}
		sigPH, err := sym.ConvertFromString(dto.Call.SigPH)
		if err != nil {
			return nil, err
		}
		ys, err := sym.ConvertFromStrings(dto.Call.Ys)
		if err != nil {
			return nil, err
		}
		return CallSpecOld{X: x, SigPH: sigPH, Ys: ys}, nil
	case Fwd:
		x, err := sym.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		y, err := sym.ConvertFromString(dto.Fwd.Y)
		if err != nil {
			return nil, err
		}
		return FwdSpec{X: x, Y: y}, nil
	default:
		panic(ErrUnexpectedTermKind(dto.K))
	}
}

func ErrUnexpectedTermKind(k TermKind) error {
	return fmt.Errorf("unexpected term kind: %v", k)
}

func ErrUnexpectedSemKind(k SemKind) error {
	return fmt.Errorf("unexpected sem kind: %v", k)
}

func DataFromTermRec(r TermRec) (TermRecDS, error) {
	switch rec := r.(type) {
	case CloseRec:
		return TermRecDS{
			K:     closeKind,
			Close: &closeRecDS{sym.ConvertToString(rec.X)},
		}, nil
	case WaitRec:
		dto, err := dataFromTermSpec(rec.Cont)
		if err != nil {
			return TermRecDS{}, err
		}
		return TermRecDS{
			K: waitKind,
			Wait: &waitRecDS{
				X:    sym.ConvertToString(rec.X),
				Cont: dto,
			},
		}, nil
	case SendRec:
		return TermRecDS{
			K: sendKind,
			Send: &sendRecDS{
				X: sym.ConvertToString(rec.X),
				A: id.ConvertToString(rec.A),
				B: id.ConvertToString(rec.B),
			},
		}, nil
	case RecvRec:
		dto, err := dataFromTermSpec(rec.Cont)
		if err != nil {
			return TermRecDS{}, err
		}
		return TermRecDS{
			K: recvKind,
			Recv: &recvRecDS{
				X:    sym.ConvertToString(rec.X),
				Y:    sym.ConvertToString(rec.Y),
				Cont: dto,
			},
		}, nil
	case LabRec:
		return TermRecDS{
			K:   labKind,
			Lab: &labRecDS{sym.ConvertToString(rec.X), string(rec.Label)},
		}, nil
	case CaseRec:
		brs := []branchRecDS{}
		for l, cont := range rec.Conts {
			dto, err := dataFromTermSpec(cont)
			if err != nil {
				return TermRecDS{}, err
			}
			brs = append(brs, branchRecDS{Label: string(l), Cont: dto})
		}
		return TermRecDS{
			K: caseKind,
			Case: &caseRecDS{
				X:        sym.ConvertToString(rec.X),
				Branches: brs,
			},
		}, nil
	case FwdRec:
		return TermRecDS{
			K: fwdKind,
			Fwd: &fwdRecDS{
				X: sym.ConvertToString(rec.X),
				B: id.ConvertToString(rec.B),
			},
		}, nil
	default:
		panic(ErrTermTypeUnexpected(rec))
	}
}

func DataToTermRec(dto TermRecDS) (TermRec, error) {
	switch dto.K {
	case closeKind:
		a, err := sym.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return CloseRec{X: a}, nil
	case waitKind:
		x, err := sym.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTermSpec(dto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitRec{X: x, Cont: cont}, nil
	case sendKind:
		x, err := sym.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		a, err := id.ConvertFromString(dto.Send.A)
		if err != nil {
			return nil, err
		}
		return SendRec{X: x, A: a}, nil
	case recvKind:
		x, err := sym.ConvertFromString(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := sym.ConvertFromString(dto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTermSpec(dto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return RecvRec{X: x, Y: y, Cont: cont}, nil
	case labKind:
		a, err := sym.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		return LabRec{X: a, Label: sym.ADT(dto.Lab.Label)}, nil
	case caseKind:
		x, err := sym.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[sym.ADT]TermSpec, len(dto.Case.Branches))
		for _, branch := range dto.Case.Branches {
			cont, err := dataToTermSpec(branch.Cont)
			if err != nil {
				return nil, err
			}
			conts[sym.ADT(branch.Label)] = cont
		}
		return CaseRec{X: x, Conts: conts}, nil
	case fwdKind:
		x, err := sym.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		b, err := id.ConvertFromString(dto.Fwd.B)
		if err != nil {
			return nil, err
		}
		return FwdRec{X: x, B: b}, nil
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func dataFromTermSpec(s TermSpec) (TermSpecDS, error) {
	switch spec := s.(type) {
	case CloseSpec:
		return TermSpecDS{
			K:     closeKind,
			Close: &closeSpecDS{sym.ConvertToString(spec.CommPH)},
		}, nil
	case WaitSpec:
		dto, err := dataFromTermSpec(spec.ContTS)
		if err != nil {
			return TermSpecDS{}, err
		}
		return TermSpecDS{
			K: waitKind,
			Wait: &waitSpecDS{
				X:    sym.ConvertToString(spec.CommPH),
				Cont: dto,
			},
		}, nil
	case SendSpec:
		return TermSpecDS{
			K: sendKind,
			Send: &sendSpecDS{
				X: sym.ConvertToString(spec.CommPH),
				Y: sym.ConvertToString(spec.ValPH),
			},
		}, nil
	case RecvSpec:
		dto, err := dataFromTermSpec(spec.ContTS)
		if err != nil {
			return TermSpecDS{}, err
		}
		return TermSpecDS{
			K: recvKind,
			Recv: &recvSpecDS{
				X:    sym.ConvertToString(spec.CommPH),
				Y:    sym.ConvertToString(spec.BindPH),
				Cont: dto,
			},
		}, nil
	case LabSpec:
		return TermSpecDS{
			K:   labKind,
			Lab: &labSpecDS{sym.ConvertToString(spec.CommPH), string(spec.Label)},
		}, nil
	case CaseSpec:
		brs := []branchSpecDS{}
		for l, cont := range spec.Conts {
			dto, err := dataFromTermSpec(cont)
			if err != nil {
				return TermSpecDS{}, err
			}
			brs = append(brs, branchSpecDS{Label: string(l), Cont: dto})
		}
		return TermSpecDS{
			K: caseKind,
			Case: &caseSpecDS{
				X:        sym.ConvertToString(spec.CommPH),
				Branches: brs,
			},
		}, nil
	case FwdSpec:
		return TermSpecDS{
			K: fwdKind,
			Fwd: &fwdSpecDS{
				X: sym.ConvertToString(spec.X),
				Y: sym.ConvertToString(spec.Y),
			},
		}, nil
	default:
		panic(ErrTermTypeUnexpected(spec))
	}
}

func dataToTermSpec(dto TermSpecDS) (TermSpec, error) {
	switch dto.K {
	case closeKind:
		a, err := sym.ConvertFromString(dto.Close.X)
		if err != nil {
			return nil, err
		}
		return CloseSpec{CommPH: a}, nil
	case waitKind:
		x, err := sym.ConvertFromString(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTermSpec(dto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitSpec{CommPH: x, ContTS: cont}, nil
	case sendKind:
		x, err := sym.ConvertFromString(dto.Send.X)
		if err != nil {
			return nil, err
		}
		y, err := sym.ConvertFromString(dto.Send.Y)
		if err != nil {
			return nil, err
		}
		return SendSpec{CommPH: x, ValPH: y}, nil
	case recvKind:
		x, err := sym.ConvertFromString(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := sym.ConvertFromString(dto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTermSpec(dto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return RecvSpec{CommPH: x, BindPH: y, ContTS: cont}, nil
	case labKind:
		x, err := sym.ConvertFromString(dto.Lab.X)
		if err != nil {
			return nil, err
		}
		return LabSpec{CommPH: x, Label: sym.ADT(dto.Lab.Label)}, nil
	case caseKind:
		x, err := sym.ConvertFromString(dto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[sym.ADT]TermSpec, len(dto.Case.Branches))
		for _, b := range dto.Case.Branches {
			cont, err := dataToTermSpec(b.Cont)
			if err != nil {
				return nil, err
			}
			conts[sym.ADT(b.Label)] = cont
		}
		return CaseSpec{CommPH: x, Conts: conts}, nil
	case fwdKind:
		x, err := sym.ConvertFromString(dto.Fwd.X)
		if err != nil {
			return nil, err
		}
		y, err := sym.ConvertFromString(dto.Fwd.Y)
		if err != nil {
			return nil, err
		}
		return FwdSpec{X: x, Y: y}, nil
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func errUnexpectedTermKind(k termKind) error {
	return fmt.Errorf("unexpected term kind: %v", k)
}
