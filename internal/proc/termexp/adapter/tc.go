package adapter

import (
	"orglang/go-engine/internal/adt/contkey"
	"fmt"
	"sort"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	termexp "orglang/go-engine/internal/proc/termexp/core"
)

func DataFromExpRec(r termexp.CompRec) (ExpRec, error) {
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

func DataToExpRec(dto ExpRec) (termexp.CompRec, error) {
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
		conts := make(map[uniqsym.ADT]termexp.CompSpec, len(dto.Case.Branches))
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

func dataFromExpSpec(s termexp.CompSpec) (ExpSpec, error) {
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

func dataToExpSpec(dto ExpSpec) (termexp.CompSpec, error) {
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
		conts := make(map[uniqsym.ADT]termexp.CompSpec, len(dto.Case.Branches))
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


func dataFromDescRec(rec termexp.DescRec) (DescData,error) {
	if rec==nil { return DescData{},fmt.Errorf("nil process Description record") }
	dto:=DescData{ExpCK:contkey.ConvertToInt(rec.ContentKey())}
	statesFromDescRec(0,rec,&dto)
	return dto,nil
}
func statesFromDescRec(parent int64, rec termexp.DescRec, dto *DescData) int64 {
	key:=contkey.ConvertToInt(rec.ContentKey())
	var kind int16
	d:=StateSpec{}
	child:=func(r termexp.DescRec) int64 { return statesFromDescRec(key,r,dto) }
	switch r:=rec.(type) {
	case termexp.DescCloseRec: kind=termexp.CloseKind; d.ContChnlPH=string(r.ContChnlPH)
	case termexp.DescWaitRec: kind=termexp.WaitKind; d.ContChnlPH=string(r.ContChnlPH); d.ContExpCK=child(r.ContExp)
	case termexp.DescSendRec: kind=termexp.SendKind; d.ConnChnlPH=string(r.ConnChnlPH); d.ValChnlPH=string(r.ValChnlPH)
	case termexp.DescRecvRec: kind=termexp.RecvKind; d.ConnChnlPH=string(r.ConnChnlPH); d.NewChnlPH=string(r.NewChnlPH); d.ContExpCK=child(r.ContExp)
	case termexp.DescLabRec: kind=termexp.LabKind; d.ConnChnlPH=string(r.ConnChnlPH); d.ValLabQN=uniqsym.ConvertToString(r.ValLabQN); d.ContExpCK=child(r.ContExp)
	case termexp.DescCaseRec:
		kind=termexp.CaseKind; d.ConnChnlPH=string(r.ConnChnlPH); d.ContExpCKs=make(map[string]int64,len(r.ContExps))
		labels:=make([]uniqsym.ADT,0,len(r.ContExps)); for q:=range r.ContExps { labels=append(labels,q) }
		sort.Slice(labels,func(i,j int)bool{return uniqsym.ConvertToString(labels[i])<uniqsym.ConvertToString(labels[j])})
		for _,q:=range labels { d.ContExpCKs[uniqsym.ConvertToString(q)]=child(r.ContExps[q]) }
	case termexp.DescLinkRec:
		kind=termexp.LinkKind; d.ProcTermQN=uniqsym.ConvertToString(r.ProcTermQN); d.X=identity.ConvertToString(r.X); d.Ys=make([]string,len(r.Ys)); for i,v:=range r.Ys { d.Ys[i]=identity.ConvertToString(v) }
	case termexp.DescFwdRec: kind=termexp.FwdKind; d.ConnChnlPH=string(r.ConnChnlPH); d.ContChnlPH=string(r.ContChnlPH)
	case termexp.DescCallRec:
		kind=termexp.CallKind; d.ConnChnlPH=string(r.ConnChnlPH); d.NewChnlPH=string(r.NewChnlPH); d.ProcTermQN=uniqsym.ConvertToString(r.ProcTermQN); d.ValChnlPHs=make([]string,len(r.ValChnlPHs)); for i,v:=range r.ValChnlPHs { d.ValChnlPHs[i]=string(v) }; d.ContExpCK=child(r.ContExp)
	case termexp.DescSpawnRec:
		kind=termexp.SpawnKind; d.ConnChnlPH=string(r.ConnChnlPH); d.ProcTermQN=uniqsym.ConvertToString(r.ProcTermQN); d.NewChnlPHs=make([]string,len(r.NewChnlPHs)); for i,v:=range r.NewChnlPHs { d.NewChnlPHs[i]=string(v) }; d.ContExpCK=child(r.ContExp)
	case termexp.DescAcquireRec: kind=termexp.AcquireKind; d.ConnChnlPH=string(r.ConnChnlPH); d.ContExpCK=child(r.ContExp)
	case termexp.DescAcceptRec: kind=termexp.AcceptKind; d.ConnChnlPH=string(r.ConnChnlPH); d.ContExpCK=child(r.ContExp)
	case termexp.DescDetachRec: kind=termexp.DetachKind; d.ConnChnlPH=string(r.ConnChnlPH)
	case termexp.DescReleaseRec: kind=termexp.ReleaseKind; d.ConnChnlPH=string(r.ConnChnlPH)
	default: panic(fmt.Sprintf("unexpected process Description record %T",rec))
	}
	dto.States=append(dto.States,DescState{ExpCK:key,SupExpCK:parent,K:kind,Spec:d})
	return key
}

func dataToDescRec(dto DescData) (termexp.DescRec,error) {
	states:=make(map[int64]DescState,len(dto.States)); for _,s:=range dto.States { states[s.ExpCK]=s }
	root,ok:=states[dto.ExpCK]; if !ok { return nil,fmt.Errorf("process Description root %d is missing",dto.ExpCK) }
	return statesToDescRec(states,root)
}
func statesToDescRec(states map[int64]DescState, st DescState) (termexp.DescRec,error) {
	key,err:=contkey.ConvertFromInt(st.ExpCK); if err!=nil{return nil,err}
	d:=st.Spec
	child:=func(ck int64)(termexp.DescRec,error){ s,ok:=states[ck]; if !ok{return nil,fmt.Errorf("process Description child %d is missing",ck)}; if s.SupExpCK!=st.ExpCK{return nil,fmt.Errorf("process Description child %d has unexpected parent",ck)}; return statesToDescRec(states,s) }
	ph:=func(v string)(symbol.ADT,error){return symbol.ConvertFromString(v)}
	switch st.K {
	case termexp.CloseKind: p,e:=ph(d.ContChnlPH); return termexp.DescCloseRec{ExpCK:key,ContChnlPH:p},e
	case termexp.WaitKind: p,e:=ph(d.ContChnlPH); if e!=nil{return nil,e}; c,e:=child(d.ContExpCK); return termexp.DescWaitRec{ExpCK:key,ContChnlPH:p,ContExp:c},e
	case termexp.SendKind: x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; y,e:=ph(d.ValChnlPH); return termexp.DescSendRec{ExpCK:key,ConnChnlPH:x,ValChnlPH:y},e
	case termexp.RecvKind: x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; y,e:=ph(d.NewChnlPH); if e!=nil{return nil,e}; c,e:=child(d.ContExpCK); return termexp.DescRecvRec{ExpCK:key,ConnChnlPH:x,NewChnlPH:y,ContExp:c},e
	case termexp.LabKind: x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; q,e:=uniqsym.ConvertFromString(d.ValLabQN); if e!=nil{return nil,e}; c,e:=child(d.ContExpCK); return termexp.DescLabRec{ExpCK:key,ConnChnlPH:x,ValLabQN:q,ContExp:c},e
	case termexp.CaseKind:
		x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; cs:=make(map[uniqsym.ADT]termexp.DescRec,len(d.ContExpCKs)); for name,ck:=range d.ContExpCKs { q,e:=uniqsym.ConvertFromString(name); if e!=nil{return nil,e}; c,e:=child(ck); if e!=nil{return nil,e}; cs[q]=c }; return termexp.DescCaseRec{ExpCK:key,ConnChnlPH:x,ContExps:cs},nil
	case termexp.LinkKind:
		q,e:=uniqsym.ConvertFromString(d.ProcTermQN); if e!=nil{return nil,e}; x,e:=identity.ConvertFromString(d.X); if e!=nil{return nil,e}; ys:=make([]identity.ADT,len(d.Ys)); for i,v:=range d.Ys { ys[i],e=identity.ConvertFromString(v); if e!=nil{return nil,e} }; return termexp.DescLinkRec{ExpCK:key,ProcTermQN:q,X:x,Ys:ys},nil
	case termexp.FwdKind: x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; y,e:=ph(d.ContChnlPH); return termexp.DescFwdRec{ExpCK:key,ConnChnlPH:x,ContChnlPH:y},e
	case termexp.CallKind:
		x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; y,e:=ph(d.NewChnlPH); if e!=nil{return nil,e}; q,e:=uniqsym.ConvertFromString(d.ProcTermQN); if e!=nil{return nil,e}; chs:=make([]symbol.ADT,len(d.ValChnlPHs)); for i,v:=range d.ValChnlPHs { chs[i],e=ph(v); if e!=nil{return nil,e} }; c,e:=child(d.ContExpCK); return termexp.DescCallRec{ExpCK:key,ConnChnlPH:x,NewChnlPH:y,ProcTermQN:q,ValChnlPHs:chs,ContExp:c},e
	case termexp.SpawnKind:
		x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; q,e:=uniqsym.ConvertFromString(d.ProcTermQN); if e!=nil{return nil,e}; chs:=make([]symbol.ADT,len(d.NewChnlPHs)); for i,v:=range d.NewChnlPHs { chs[i],e=ph(v); if e!=nil{return nil,e} }; c,e:=child(d.ContExpCK); return termexp.DescSpawnRec{ExpCK:key,ConnChnlPH:x,ProcTermQN:q,NewChnlPHs:chs,ContExp:c},e
	case termexp.AcquireKind: x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; c,e:=child(d.ContExpCK); return termexp.DescAcquireRec{ExpCK:key,ConnChnlPH:x,ContExp:c},e
	case termexp.AcceptKind: x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; c,e:=child(d.ContExpCK); return termexp.DescAcceptRec{ExpCK:key,ConnChnlPH:x,ContExp:c},e
	case termexp.DetachKind: x,e:=ph(d.ConnChnlPH); return termexp.DescDetachRec{ExpCK:key,ConnChnlPH:x},e
	case termexp.ReleaseKind: x,e:=ph(d.ConnChnlPH); return termexp.DescReleaseRec{ExpCK:key,ConnChnlPH:x},e
	default: return nil,fmt.Errorf("unexpected process Description kind %d",st.K)
	}
}
