package adapter

import (
	"encoding/json"
	"fmt"
	"sort"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

type DescData struct {
	ExpCK int64 `db:"exp_ck"`
	States []DescState `db:"states"`
}
type DescState struct {
	ExpCK int64 `db:"exp_ck"`
	SupExpCK int64 `db:"sup_exp_ck"`
	K int16 `db:"kind"`
	Spec json.RawMessage `db:"spec"`
}
type descSpecData struct {
	ContChnlPH string `json:"cont_chnl_ph,omitempty"`
	ConnChnlPH string `json:"conn_chnl_ph,omitempty"`
	ValChnlPH string `json:"val_chnl_ph,omitempty"`
	NewChnlPH string `json:"new_chnl_ph,omitempty"`
	ContExpCK int64 `json:"cont_exp_ck,omitempty"`
	ContExpCKs map[string]int64 `json:"cont_exp_cks,omitempty"`
	ProcTermQN string `json:"proc_term_qn,omitempty"`
	X string `json:"x,omitempty"`
	Ys []string `json:"ys,omitempty"`
	ValChnlPHs []string `json:"val_chnl_phs,omitempty"`
	NewChnlPHs []string `json:"new_chnl_phs,omitempty"`
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
	d:=descSpecData{}
	child:=func(r termexp.DescRec) int64 { return statesFromDescRec(key,r,dto) }
	switch r:=rec.(type) {
	case termexp.DescCloseRec: kind=termexp.CloseKind; d.ContChnlPH=string(r.ContChnlPH)
	case termexp.DescWaitRec: kind=termexp.WaitKind; d.ContChnlPH=string(r.ContChnlPH); d.ContExpCK=child(r.ContExp)
	case termexp.DescSendRec: kind=termexp.SendKind; d.ConnChnlPH=string(r.ConnChnlPH); d.ValChnlPH=string(r.ValChnlPH)
	case termexp.DescRecvRec: kind=termexp.RecvKind; d.ConnChnlPH=string(r.ConnChnlPH); d.NewChnlPH=string(r.NewChnlPH); d.ContExpCK=child(r.ContExp)
	case termexp.DescLabRec: kind=termexp.LabKind; d.ConnChnlPH=string(r.ConnChnlPH); d.ProcTermQN=uniqsym.ConvertToString(r.ValLabQN); d.ContExpCK=child(r.ContExp)
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
	raw,err:=json.Marshal(d); if err!=nil { panic(err) }
	dto.States=append(dto.States,DescState{ExpCK:key,SupExpCK:parent,K:kind,Spec:raw})
	return key
}

func dataToDescRec(dto DescData) (termexp.DescRec,error) {
	states:=make(map[int64]DescState,len(dto.States)); for _,s:=range dto.States { states[s.ExpCK]=s }
	root,ok:=states[dto.ExpCK]; if !ok { return nil,fmt.Errorf("process Description root %d is missing",dto.ExpCK) }
	return statesToDescRec(states,root)
}
func statesToDescRec(states map[int64]DescState, st DescState) (termexp.DescRec,error) {
	key,err:=contkey.ConvertFromInt(st.ExpCK); if err!=nil{return nil,err}
	var d descSpecData; if err=json.Unmarshal(st.Spec,&d); err!=nil{return nil,err}
	child:=func(ck int64)(termexp.DescRec,error){ s,ok:=states[ck]; if !ok{return nil,fmt.Errorf("process Description child %d is missing",ck)}; if s.SupExpCK!=st.ExpCK{return nil,fmt.Errorf("process Description child %d has unexpected parent",ck)}; return statesToDescRec(states,s) }
	ph:=func(v string)(symbol.ADT,error){return symbol.ConvertFromString(v)}
	switch st.K {
	case termexp.CloseKind: p,e:=ph(d.ContChnlPH); return termexp.DescCloseRec{ExpCK:key,ContChnlPH:p},e
	case termexp.WaitKind: p,e:=ph(d.ContChnlPH); if e!=nil{return nil,e}; c,e:=child(d.ContExpCK); return termexp.DescWaitRec{ExpCK:key,ContChnlPH:p,ContExp:c},e
	case termexp.SendKind: x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; y,e:=ph(d.ValChnlPH); return termexp.DescSendRec{ExpCK:key,ConnChnlPH:x,ValChnlPH:y},e
	case termexp.RecvKind: x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; y,e:=ph(d.NewChnlPH); if e!=nil{return nil,e}; c,e:=child(d.ContExpCK); return termexp.DescRecvRec{ExpCK:key,ConnChnlPH:x,NewChnlPH:y,ContExp:c},e
	case termexp.LabKind: x,e:=ph(d.ConnChnlPH); if e!=nil{return nil,e}; q,e:=uniqsym.ConvertFromString(d.ProcTermQN); if e!=nil{return nil,e}; c,e:=child(d.ContExpCK); return termexp.DescLabRec{ExpCK:key,ConnChnlPH:x,ValLabQN:q,ContExp:c},e
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
