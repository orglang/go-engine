package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

func TestContentKeyIncludesLabContinuation(t *testing.T) {
	label := uniqsym.New(symbol.New("accepted"))
	base := termexp.LabSpec{
		ConnChnlPH: symbol.New("conn"),
		ValLabQN:   label,
		ContExp:    termexp.CloseSpec{ContChnlPH: symbol.New("left")},
	}
	changed := termexp.LabSpec{
		ConnChnlPH: base.ConnChnlPH,
		ValLabQN:   label,
		ContExp:    termexp.CloseSpec{ContChnlPH: symbol.New("right")},
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
	left := termexp.CaseSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExps: map[uniqsym.ADT]termexp.ExpSpec{
			leftLabel:  termexp.CloseSpec{ContChnlPH: symbol.New("left")},
			rightLabel: termexp.CloseSpec{ContChnlPH: symbol.New("right")},
		},
	}
	right := termexp.CaseSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExps: map[uniqsym.ADT]termexp.ExpSpec{
			rightLabel: termexp.CloseSpec{ContChnlPH: symbol.New("right")},
			leftLabel:  termexp.CloseSpec{ContChnlPH: symbol.New("left")},
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
