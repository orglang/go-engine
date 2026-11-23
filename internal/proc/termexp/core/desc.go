package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

// DescSpec is the declarative expression stored as a Term Definition.
type DescSpec interface {
	descSpec()
}

type DescCloseSpec struct { ContChnlPH symbol.ADT }
func (DescCloseSpec) descSpec() {}
type DescWaitSpec struct { ContChnlPH symbol.ADT; ContExp DescSpec }
func (DescWaitSpec) descSpec() {}
type DescSendSpec struct { ConnChnlPH symbol.ADT; ValChnlPH symbol.ADT }
func (DescSendSpec) descSpec() {}
type DescRecvSpec struct { ConnChnlPH symbol.ADT; NewChnlPH symbol.ADT; ContExp DescSpec }
func (DescRecvSpec) descSpec() {}
type DescLabSpec struct { ConnChnlPH symbol.ADT; ValLabQN uniqsym.ADT; ContExp DescSpec }
func (DescLabSpec) descSpec() {}
type DescCaseSpec struct { ConnChnlPH symbol.ADT; ContExps map[uniqsym.ADT]DescSpec }
func (DescCaseSpec) descSpec() {}
type DescLinkSpec struct { ProcTermQN uniqsym.ADT; X identity.ADT; Ys []identity.ADT }
func (DescLinkSpec) descSpec() {}
type DescFwdSpec struct { ConnChnlPH symbol.ADT; ContChnlPH symbol.ADT }
func (DescFwdSpec) descSpec() {}
type DescCallSpec struct { ConnChnlPH symbol.ADT; NewChnlPH symbol.ADT; ProcTermQN uniqsym.ADT; ValChnlPHs []symbol.ADT; ContExp DescSpec }
func (DescCallSpec) descSpec() {}
type DescSpawnSpec struct { ConnChnlPH symbol.ADT; ProcTermQN uniqsym.ADT; NewChnlPHs []symbol.ADT; ContExp DescSpec }
func (DescSpawnSpec) descSpec() {}
type DescAcquireSpec struct { ConnChnlPH symbol.ADT; ContExp DescSpec }
func (DescAcquireSpec) descSpec() {}
type DescAcceptSpec struct { ConnChnlPH symbol.ADT; ContExp DescSpec }
func (DescAcceptSpec) descSpec() {}
type DescDetachSpec struct { ConnChnlPH symbol.ADT }
func (DescDetachSpec) descSpec() {}
type DescReleaseSpec struct { ConnChnlPH symbol.ADT }
func (DescReleaseSpec) descSpec() {}

// DescRec is the recursive, content-addressed domain representation of a Definition.
type DescRec interface { descRec(); ContentKey() contkey.ADT }
type DescCloseRec struct { ExpCK contkey.ADT; ContChnlPH symbol.ADT }
func (DescCloseRec) descRec() {}
func (r DescCloseRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescWaitRec struct { ExpCK contkey.ADT; ContChnlPH symbol.ADT; ContExp DescRec }
func (DescWaitRec) descRec() {}
func (r DescWaitRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescSendRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ValChnlPH symbol.ADT }
func (DescSendRec) descRec() {}
func (r DescSendRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescRecvRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; NewChnlPH symbol.ADT; ContExp DescRec }
func (DescRecvRec) descRec() {}
func (r DescRecvRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescLabRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ValLabQN uniqsym.ADT; ContExp DescRec }
func (DescLabRec) descRec() {}
func (r DescLabRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescCaseRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ContExps map[uniqsym.ADT]DescRec }
func (DescCaseRec) descRec() {}
func (r DescCaseRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescLinkRec struct { ExpCK contkey.ADT; ProcTermQN uniqsym.ADT; X identity.ADT; Ys []identity.ADT }
func (DescLinkRec) descRec() {}
func (r DescLinkRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescFwdRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ContChnlPH symbol.ADT }
func (DescFwdRec) descRec() {}
func (r DescFwdRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescCallRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; NewChnlPH symbol.ADT; ProcTermQN uniqsym.ADT; ValChnlPHs []symbol.ADT; ContExp DescRec }
func (DescCallRec) descRec() {}
func (r DescCallRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescSpawnRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ProcTermQN uniqsym.ADT; NewChnlPHs []symbol.ADT; ContExp DescRec }
func (DescSpawnRec) descRec() {}
func (r DescSpawnRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescAcquireRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ContExp DescRec }
func (DescAcquireRec) descRec() {}
func (r DescAcquireRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescAcceptRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ContExp DescRec }
func (DescAcceptRec) descRec() {}
func (r DescAcceptRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescDetachRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT }
func (DescDetachRec) descRec() {}
func (r DescDetachRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescReleaseRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT }
func (DescReleaseRec) descRec() {}
func (r DescReleaseRec) ContentKey() contkey.ADT { return r.ExpCK }

// ConvertSpecToRec converts a declarative Definition into its content-addressed domain record.
func ConvertSpecToRec(spec DescSpec) (DescRec, error) {
	key, _, err := DescContentKey(spec)
	if err != nil { return nil, err }
	switch s := spec.(type) {
	case DescCloseSpec: return DescCloseRec{key,s.ContChnlPH}, nil
	case DescWaitSpec:
		c, e := ConvertSpecToRec(s.ContExp); if e != nil { return nil,e }; return DescWaitRec{key,s.ContChnlPH,c},nil
	case DescSendSpec: return DescSendRec{key,s.ConnChnlPH,s.ValChnlPH},nil
	case DescRecvSpec:
		c,e:=ConvertSpecToRec(s.ContExp); if e!=nil{return nil,e}; return DescRecvRec{key,s.ConnChnlPH,s.NewChnlPH,c},nil
	case DescLabSpec:
		c,e:=ConvertSpecToRec(s.ContExp); if e!=nil{return nil,e}; return DescLabRec{key,s.ConnChnlPH,s.ValLabQN,c},nil
	case DescCaseSpec:
		cs:=make(map[uniqsym.ADT]DescRec,len(s.ContExps)); for q,v:=range s.ContExps { c,e:=ConvertSpecToRec(v); if e!=nil{return nil,e}; cs[q]=c }; return DescCaseRec{key,s.ConnChnlPH,cs},nil
	case DescLinkSpec: return DescLinkRec{key,s.ProcTermQN,s.X,s.Ys},nil
	case DescFwdSpec: return DescFwdRec{key,s.ConnChnlPH,s.ContChnlPH},nil
	case DescCallSpec:
		c,e:=ConvertSpecToRec(s.ContExp); if e!=nil{return nil,e}; return DescCallRec{key,s.ConnChnlPH,s.NewChnlPH,s.ProcTermQN,s.ValChnlPHs,c},nil
	case DescSpawnSpec:
		c,e:=ConvertSpecToRec(s.ContExp); if e!=nil{return nil,e}; return DescSpawnRec{key,s.ConnChnlPH,s.ProcTermQN,s.NewChnlPHs,c},nil
	case DescAcquireSpec:
		c,e:=ConvertSpecToRec(s.ContExp); if e!=nil{return nil,e}; return DescAcquireRec{key,s.ConnChnlPH,c},nil
	case DescAcceptSpec:
		c,e:=ConvertSpecToRec(s.ContExp); if e!=nil{return nil,e}; return DescAcceptRec{key,s.ConnChnlPH,c},nil
	case DescDetachSpec: return DescDetachRec{key,s.ConnChnlPH},nil
	case DescReleaseSpec: return DescReleaseRec{key,s.ConnChnlPH},nil
	default: return nil,fmt.Errorf("unexpected process Description spec %T",spec)
	}
}

func ConvertRecToSpec(rec DescRec) DescSpec {
	switch r:=rec.(type) {
	case DescCloseRec: return DescCloseSpec{r.ContChnlPH}
	case DescWaitRec: return DescWaitSpec{r.ContChnlPH,ConvertRecToSpec(r.ContExp)}
	case DescSendRec: return DescSendSpec{r.ConnChnlPH,r.ValChnlPH}
	case DescRecvRec: return DescRecvSpec{r.ConnChnlPH,r.NewChnlPH,ConvertRecToSpec(r.ContExp)}
	case DescLabRec: return DescLabSpec{r.ConnChnlPH,r.ValLabQN,ConvertRecToSpec(r.ContExp)}
	case DescCaseRec:
		cs:=make(map[uniqsym.ADT]DescSpec,len(r.ContExps)); for q,v:=range r.ContExps { cs[q]=ConvertRecToSpec(v) }; return DescCaseSpec{r.ConnChnlPH,cs}
	case DescLinkRec: return DescLinkSpec{r.ProcTermQN,r.X,r.Ys}
	case DescFwdRec: return DescFwdSpec{r.ConnChnlPH,r.ContChnlPH}
	case DescCallRec: return DescCallSpec{r.ConnChnlPH,r.NewChnlPH,r.ProcTermQN,r.ValChnlPHs,ConvertRecToSpec(r.ContExp)}
	case DescSpawnRec: return DescSpawnSpec{r.ConnChnlPH,r.ProcTermQN,r.NewChnlPHs,ConvertRecToSpec(r.ContExp)}
	case DescAcquireRec: return DescAcquireSpec{r.ConnChnlPH,ConvertRecToSpec(r.ContExp)}
	case DescAcceptRec: return DescAcceptSpec{r.ConnChnlPH,ConvertRecToSpec(r.ContExp)}
	case DescDetachRec: return DescDetachSpec{r.ConnChnlPH}
	case DescReleaseRec: return DescReleaseSpec{r.ConnChnlPH}
	default: panic(fmt.Sprintf("unexpected process Description rec %T",rec))
	}
}
