package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

// DescSpec is the declarative expression stored as a Term Definition.
type DescSpec interface { descSpec() }
type DescAcquireSpec struct { ConnChnlPH symbol.ADT; ContExp DescSpec }
func (DescAcquireSpec) descSpec() {}
type DescAcceptSpec struct { ConnChnlPH symbol.ADT; ContExp DescSpec }
func (DescAcceptSpec) descSpec() {}
type DescHireSpec struct { ConnChnlPH symbol.ADT; ProcTermQN uniqsym.ADT; ContExp DescSpec }
func (DescHireSpec) descSpec() {}
type DescApplySpec struct { ConnChnlPH symbol.ADT; ProcTermQN uniqsym.ADT; ContExp DescSpec }
func (DescApplySpec) descSpec() {}
type DescReleaseSpec struct { ConnChnlPH symbol.ADT }
func (DescReleaseSpec) descSpec() {}
type DescDetachSpec struct { ConnChnlPH symbol.ADT }
func (DescDetachSpec) descSpec() {}
type DescFireSpec struct { ConnChnlPH symbol.ADT; ProcTermQN uniqsym.ADT }
func (DescFireSpec) descSpec() {}
type DescQuitSpec struct { ConnChnlPH symbol.ADT; ProcTermQN uniqsym.ADT }
func (DescQuitSpec) descSpec() {}
type DescSpawnSpec struct { ProcTermQN uniqsym.ADT; ProcCompQNs []uniqsym.ADT }
func (DescSpawnSpec) descSpec() {}

// DescRec is the recursive, content-addressed domain representation of a Definition.
type DescRec interface { descRec(); ContentKey() contkey.ADT }
type DescAcquireRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ContExp DescRec }
func (DescAcquireRec) descRec() {}
func (r DescAcquireRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescAcceptRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ContExp DescRec }
func (DescAcceptRec) descRec() {}
func (r DescAcceptRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescHireRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ProcTermQN uniqsym.ADT; ContExp DescRec }
func (DescHireRec) descRec() {}
func (r DescHireRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescApplyRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ProcTermQN uniqsym.ADT; ContExp DescRec }
func (DescApplyRec) descRec() {}
func (r DescApplyRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescReleaseRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT }
func (DescReleaseRec) descRec() {}
func (r DescReleaseRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescDetachRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT }
func (DescDetachRec) descRec() {}
func (r DescDetachRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescFireRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ProcTermQN uniqsym.ADT }
func (DescFireRec) descRec() {}
func (r DescFireRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescQuitRec struct { ExpCK contkey.ADT; ConnChnlPH symbol.ADT; ProcTermQN uniqsym.ADT }
func (DescQuitRec) descRec() {}
func (r DescQuitRec) ContentKey() contkey.ADT { return r.ExpCK }
type DescSpawnRec struct { ExpCK contkey.ADT; ProcTermQN uniqsym.ADT; ProcCompQNs []uniqsym.ADT }
func (DescSpawnRec) descRec() {}
func (r DescSpawnRec) ContentKey() contkey.ADT { return r.ExpCK }

func ConvertSpecToRec(spec DescSpec) (DescRec,error) {
	key,_,err:=DescContentKey(spec); if err!=nil{return nil,err}
	switch s:=spec.(type) {
	case DescAcquireSpec: c,e:=ConvertSpecToRec(s.ContExp); if e!=nil{return nil,e}; return DescAcquireRec{key,s.ConnChnlPH,c},nil
	case DescAcceptSpec: c,e:=ConvertSpecToRec(s.ContExp); if e!=nil{return nil,e}; return DescAcceptRec{key,s.ConnChnlPH,c},nil
	case DescHireSpec: c,e:=ConvertSpecToRec(s.ContExp); if e!=nil{return nil,e}; return DescHireRec{key,s.ConnChnlPH,s.ProcTermQN,c},nil
	case DescApplySpec: c,e:=ConvertSpecToRec(s.ContExp); if e!=nil{return nil,e}; return DescApplyRec{key,s.ConnChnlPH,s.ProcTermQN,c},nil
	case DescReleaseSpec: return DescReleaseRec{key,s.ConnChnlPH},nil
	case DescDetachSpec: return DescDetachRec{key,s.ConnChnlPH},nil
	case DescFireSpec: return DescFireRec{key,s.ConnChnlPH,s.ProcTermQN},nil
	case DescQuitSpec: return DescQuitRec{key,s.ConnChnlPH,s.ProcTermQN},nil
	case DescSpawnSpec: return DescSpawnRec{key,s.ProcTermQN,s.ProcCompQNs},nil
	default: return nil,fmt.Errorf("unexpected pool Description spec %T",spec)
	}
}
func ConvertRecToSpec(rec DescRec) DescSpec {
	switch r:=rec.(type) {
	case DescAcquireRec: return DescAcquireSpec{r.ConnChnlPH,ConvertRecToSpec(r.ContExp)}
	case DescAcceptRec: return DescAcceptSpec{r.ConnChnlPH,ConvertRecToSpec(r.ContExp)}
	case DescHireRec: return DescHireSpec{r.ConnChnlPH,r.ProcTermQN,ConvertRecToSpec(r.ContExp)}
	case DescApplyRec: return DescApplySpec{r.ConnChnlPH,r.ProcTermQN,ConvertRecToSpec(r.ContExp)}
	case DescReleaseRec: return DescReleaseSpec{r.ConnChnlPH}
	case DescDetachRec: return DescDetachSpec{r.ConnChnlPH}
	case DescFireRec: return DescFireSpec{r.ConnChnlPH,r.ProcTermQN}
	case DescQuitRec: return DescQuitSpec{r.ConnChnlPH,r.ProcTermQN}
	case DescSpawnRec: return DescSpawnSpec{r.ProcTermQN,r.ProcCompQNs}
	default: panic(fmt.Sprintf("unexpected pool Description rec %T",rec))
	}
}
