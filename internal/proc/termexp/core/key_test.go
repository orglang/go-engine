package core

import (
	"testing"

	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

func TestContentKeyIncludesLabContinuation(t *testing.T) {
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

	baseKey, _, err := ContentKey(base)
	if err != nil {
		t.Fatal(err)
	}
	changedKey, _, err := ContentKey(changed)
	if err != nil {
		t.Fatal(err)
	}
	if baseKey == changedKey {
		t.Fatal("different label continuations produced the same content key")
	}
}

func TestContentKeyIsStableForCaseMapInsertionOrder(t *testing.T) {
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

	leftKey, leftJSON, err := ContentKey(left)
	if err != nil {
		t.Fatal(err)
	}
	rightKey, rightJSON, err := ContentKey(right)
	if err != nil {
		t.Fatal(err)
	}
	if leftKey != rightKey || string(leftJSON) != string(rightJSON) {
		t.Fatal("equivalent case maps produced different canonical content")
	}
}
