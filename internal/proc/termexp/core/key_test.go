package core

import (
	"testing"
	"strings"


	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

func TestCompContentKeyIncludesLabContinuation(t *testing.T) {
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

	baseKey, _, err := CompContentKey(base)
	if err != nil {
		t.Fatal(err)
	}
	changedKey, _, err := CompContentKey(changed)
	if err != nil {
		t.Fatal(err)
	}
	if baseKey == changedKey {
		t.Fatal("different label continuations produced the same content key")
	}
}

func TestCompContentKeyIsStableForCaseMapInsertionOrder(t *testing.T) {
	leftLabel := uniqsym.New(symbol.New("left"))
	rightLabel := uniqsym.New(symbol.New("right"))
	left := CaseSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExps: map[uniqsym.ADT]CompSpec{
			leftLabel:  CloseSpec{ContChnlPH: symbol.New("left")},
			rightLabel: CloseSpec{ContChnlPH: symbol.New("right")},
		},
	}
	right := CaseSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExps: map[uniqsym.ADT]CompSpec{
			rightLabel: CloseSpec{ContChnlPH: symbol.New("right")},
			leftLabel:  CloseSpec{ContChnlPH: symbol.New("left")},
		},
	}

	leftKey, leftJSON, err := CompContentKey(left)
	if err != nil {
		t.Fatal(err)
	}
	rightKey, rightJSON, err := CompContentKey(right)
	if err != nil {
		t.Fatal(err)
	}
	if leftKey != rightKey || string(leftJSON) != string(rightJSON) {
		t.Fatal("equivalent case maps produced different canonical content")
	}
}

func TestCompContentKeySerializesEveryVariant(t *testing.T) {
	cont := CloseSpec{ContChnlPH: symbol.New("next")}
	specs := []struct {
		spec CompSpec
		kind string
	}{
		{CloseSpec{}, "close"},
		{WaitSpec{ContExp: cont}, "wait"},
		{SendSpec{}, "send"},
		{RecvSpec{ContExp: cont}, "recv"},
		{LabSpec{ContExp: cont}, "lab"},
		{CaseSpec{ContExps: map[uniqsym.ADT]CompSpec{}}, "case"},
		{LinkSpec{}, "link"},
		{FwdSpec{}, "fwd"},
		{CallSpec{ContExp: cont}, "call"},
		{SpawnSpec{ContExp: cont}, "spawn"},
		{AcqureSpec{ContExp: cont}, "acquire"},
		{AcceptSpec{ContExp: cont}, "accept"},
		{DetachSpec{}, "detach"},
		{ReleaseSpec{}, "release"},
	}
	for _, tc := range specs {
		_, data, err := CompContentKey(tc.spec)
		if err != nil {
			t.Fatalf("%s: %v", tc.kind, err)
		}
		if !strings.Contains(string(data), `"kind":"`+tc.kind+`"`) {
			t.Fatalf("%s: missing discriminator in %s", tc.kind, data)
		}
	}
}

