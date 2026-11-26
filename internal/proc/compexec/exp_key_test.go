package compexec

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
		ConnChnlPH: symbol.New("conn"),
		ValLabQN:   label,
		ContExp:    CloseSpec{ContChnlPH: symbol.New("left")},
	}
	changed := LabSpec{
		ConnChnlPH: base.ConnChnlPH,
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
		ConnChnlPH: symbol.New("conn"),
		ContExps: map[uniqsym.ADT]ExpSpec{
			leftLabel:  CloseSpec{ContChnlPH: symbol.New("left")},
			rightLabel: CloseSpec{ContChnlPH: symbol.New("right")},
		},
	}
	right := CaseSpec{
		ConnChnlPH: symbol.New("conn"),
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
		{SendSpec{ConnChnlPH: symbol.New("conn"), ValChnlPH: symbol.New("value")}, "send"},
		{RecvSpec{ConnChnlPH: symbol.New("conn"), NewChnlPH: symbol.New("new"), ContExp: cont}, "recv"},
		{LabSpec{ConnChnlPH: symbol.New("conn"), ValLabQN: uniqsym.New(symbol.New("label")), ContExp: cont}, "lab"},
		{CaseSpec{ConnChnlPH: symbol.New("conn"), ContExps: map[uniqsym.ADT]ExpSpec{uniqsym.New(symbol.New("label")): cont}}, "case"},
		{LinkSpec{ProcTermQN: procTermQN, X: identity.New(), Ys: []identity.ADT{identity.New()}}, "link"},
		{FwdSpec{ConnChnlPH: symbol.New("conn"), ContChnlPH: symbol.New("next")}, "fwd"},
		{CallSpec{ConnChnlPH: symbol.New("conn"), NewChnlPH: symbol.New("new"), ProcTermQN: procTermQN, ValChnlPHs: []symbol.ADT{symbol.New("value")}, ContExp: cont}, "call"},
		{SpawnSpec{ConnChnlPH: symbol.New("conn"), ProcTermQN: procTermQN, NewChnlPHs: []symbol.ADT{symbol.New("new")}, ContExp: cont}, "spawn"},
		{AcqureSpec{ConnChnlPH: symbol.New("conn"), ContExp: cont}, "acquire"},
		{AcceptSpec{ConnChnlPH: symbol.New("conn"), ContExp: cont}, "accept"},
		{DetachSpec{ConnChnlPH: symbol.New("conn")}, "detach"},
		{ReleaseSpec{ConnChnlPH: symbol.New("conn")}, "release"},
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
