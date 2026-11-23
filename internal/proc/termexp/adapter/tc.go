package adapter

import (
	"fmt"
	"sort"

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
				Y:      symbol.ConvertToString(spec.NewChnlPH),
				ContES: dto,
			},
		}, nil
	case termexp.LabSpec:
		cont, err := dataFromExpSpec(spec.ContExp)
		if err != nil {
			return ExpSpec{}, err
		}
		return ExpSpec{
			K: LabExp,
			Lab: &LabSpec{
				X: symbol.ConvertToString(spec.ConnChnlPH),
				Label: uniqsym.ConvertToString(spec.ValLabQN),
				ContES: cont,
			},
		}, nil
	case termexp.CaseSpec:
		brs := []BranchSpec{}
		labels := make([]uniqsym.ADT, 0, len(spec.ContExps))
		for label := range spec.ContExps { labels = append(labels, label) }
		sort.Slice(labels, func(i, j int) bool { return uniqsym.ConvertToString(labels[i]) < uniqsym.ConvertToString(labels[j]) })
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
		cont, err := dataToExpSpec(dto.Lab.ContES)
		if err != nil { return nil, err }
		return termexp.LabSpec{ConnChnlPH: x, ValLabQN: label, ContExp: cont}, nil
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
			label, err := uniqsym.ConvertFromString(b.Label)
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

// canonicalExpSpec converts a domain expression into the stable JSONB storage shape.
func canonicalExpSpec(spec ExpSpec) (any, error) {
	if spec == nil {
		return nil, fmt.Errorf("nil process term expression")
	}

	base := func(kind string, fields map[string]any) map[string]any {
		fields["kind"] = kind
		return fields
	}
	nested := func(child ExpSpec) (any, error) {
		return canonicalExpSpec(child)
	}

	switch s := spec.(type) {
	case CloseSpec:
		return base("close", map[string]any{"cont_chnl_ph": string(s.ContChnlPH)}), nil
	case WaitSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("wait.cont: %w", err)
		}
		return base("wait", map[string]any{"cont_chnl_ph": string(s.ContChnlPH), "cont_exp": cont}), nil
	case SendSpec:
		return base("send", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "val_chnl_ph": string(s.ValChnlPH)}), nil
	case RecvSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("recv.cont: %w", err)
		}
		return base("recv", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "new_chnl_ph": string(s.NewChnlPH), "cont_exp": cont}), nil
	case LabSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("lab.cont: %w", err)
		}
		return base("lab", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "val_lab_qn": uniqsym.ConvertToString(s.ValLabQN), "cont_exp": cont}), nil
	case CaseSpec:
		branches := make(map[string]any, len(s.ContExps))
		for qn, child := range s.ContExps {
			name := uniqsym.ConvertToString(qn)
			if _, exists := branches[name]; exists {
				return nil, fmt.Errorf("duplicate case label %q", name)
			}
			value, err := nested(child)
			if err != nil {
				return nil, fmt.Errorf("case[%s]: %w", name, err)
			}
			branches[name] = value
		}
		return base("case", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exps": branches}), nil
	case LinkSpec:
		ys := make([]string, len(s.Ys))
		for i, id := range s.Ys {
			ys[i] = identity.ConvertToString(id)
		}
		return base("link", map[string]any{"proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "x": identity.ConvertToString(s.X), "ys": ys}), nil
	case FwdSpec:
		return base("fwd", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_chnl_ph": string(s.ContChnlPH)}), nil
	case CallSpec:
		chs := make([]string, len(s.ValChnlPHs))
		for i, ch := range s.ValChnlPHs {
			chs[i] = string(ch)
		}
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("call.cont: %w", err)
		}
		return base("call", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "new_chnl_ph": string(s.NewChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "val_chnl_phs": chs, "cont_exp": cont}), nil
	case SpawnSpec:
		chs := make([]string, len(s.NewChnlPHs))
		for i, ch := range s.NewChnlPHs {
			chs[i] = string(ch)
		}
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("spawn.cont: %w", err)
		}
		return base("spawn", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "new_chnl_phs": chs, "cont_exp": cont}), nil
	case AcqureSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("acquire.cont: %w", err)
		}
		return base("acquire", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp": cont}), nil
	case AcceptSpec:
		cont, err := nested(s.ContExp)
		if err != nil {
			return nil, fmt.Errorf("accept.cont: %w", err)
		}
		return base("accept", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp": cont}), nil
	case DetachSpec:
		return base("detach", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)}), nil
	case ReleaseSpec:
		return base("release", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)}), nil
	default:
		return nil, ErrExpTypeUnexpected(spec)
	}
}
