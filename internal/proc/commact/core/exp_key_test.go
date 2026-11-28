package core

import (
	"strings"
	"testing"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

func TestExpContentKeyIncludesLabContinuation(t *testing.T) {
	label := uniqsym.New(symbol.New("accepted"))
	base := LabSpec{
		CommChnlPH: symbol.New("conn"),
		ValLabQN:   label,
		ContExp:    CloseSpec{ContChnlPH: symbol.New("left")},
	}
	changed := LabSpec{
		CommChnlPH: base.CommChnlPH,
		ValLabQN:   label,
		ContExp:    CloseSpec{ContChnlPH: symbol.New("right")},
	}

	baseKey, _, err := ExpContentKey(base)
	if err != nil {
		t.Fatal(err)
	}
	changedKey, _, err := ExpContentKey(changed)
	if err != nil {
		t.Fatal(err)
	}
	if baseKey == changedKey {
		t.Fatal("different label continuations produced the same content key")
	}
}

func TestExpContentKeyIsStableForCaseMapInsertionOrder(t *testing.T) {
	leftLabel := uniqsym.New(symbol.New("left"))
	rightLabel := uniqsym.New(symbol.New("right"))
	left := CaseSpec{
		CommChnlPH: symbol.New("conn"),
		ContExps: map[uniqsym.ADT]ExpSpec{
			leftLabel:  CloseSpec{ContChnlPH: symbol.New("left")},
			rightLabel: CloseSpec{ContChnlPH: symbol.New("right")},
		},
	}
	right := CaseSpec{
		CommChnlPH: symbol.New("conn"),
		ContExps: map[uniqsym.ADT]ExpSpec{
			rightLabel: CloseSpec{ContChnlPH: symbol.New("right")},
			leftLabel:  CloseSpec{ContChnlPH: symbol.New("left")},
		},
	}

	leftKey, leftJSON, err := ExpContentKey(left)
	if err != nil {
		t.Fatal(err)
	}
	rightKey, rightJSON, err := ExpContentKey(right)
	if err != nil {
		t.Fatal(err)
	}
	if leftKey != rightKey || string(leftJSON) != string(rightJSON) {
		t.Fatal("equivalent case maps produced different canonical content")
	}
}

func TestExpContentKeySerializesEveryVariant(t *testing.T) {
	cont := CloseSpec{ContChnlPH: symbol.New("next")}
	procTermQN := uniqsym.New(symbol.New("process"))
	specs := []struct {
		spec ExpSpec
		kind string
	}{
		{CloseSpec{ContChnlPH: symbol.New("end")}, "close"},
		{WaitSpec{ContChnlPH: symbol.New("conn"), ContExp: cont}, "wait"},
		{SendSpec{CommChnlPH: symbol.New("conn"), ValChnlPH: symbol.New("value")}, "send"},
		{RecvSpec{CommChnlPH: symbol.New("conn"), NewChnlPH: symbol.New("new"), ContExp: cont}, "recv"},
		{LabSpec{CommChnlPH: symbol.New("conn"), ValLabQN: uniqsym.New(symbol.New("label")), ContExp: cont}, "lab"},
		{CaseSpec{CommChnlPH: symbol.New("conn"), ContExps: map[uniqsym.ADT]ExpSpec{uniqsym.New(symbol.New("label")): cont}}, "case"},
		{LinkSpec{ProcTermQN: procTermQN, X: identity.New(), Ys: []identity.ADT{identity.New()}}, "link"},
		{FwdSpec{CommChnlPH: symbol.New("conn"), ContChnlPH: symbol.New("next")}, "fwd"},
		{CallSpec{CommChnlPH: symbol.New("conn"), NewChnlPH: symbol.New("new"), ProcTermQN: procTermQN, ValChnlPHs: []symbol.ADT{symbol.New("value")}, ContExp: cont}, "call"},
		{SpawnSpec{CommChnlPH: symbol.New("conn"), ProcTermQN: procTermQN, NewChnlPHs: []symbol.ADT{symbol.New("new")}, ContExp: cont}, "spawn"},
		{AcqureSpec{CommChnlPH: symbol.New("conn"), ContExp: cont}, "acquire"},
		{AcceptSpec{CommChnlPH: symbol.New("conn"), ContExp: cont}, "accept"},
		{DetachSpec{CommChnlPH: symbol.New("conn")}, "detach"},
		{ReleaseSpec{CommChnlPH: symbol.New("conn")}, "release"},
	}
	for _, tc := range specs {
		_, data, err := ExpContentKey(tc.spec)
		if err != nil {
			t.Fatalf("%s: %v", tc.kind, err)
		}
		if !strings.Contains(string(data), `"kind":"`+tc.kind+`"`) {
			t.Fatalf("%s: missing discriminator in %s", tc.kind, data)
		}
	}
}
