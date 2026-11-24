package core

import (
	"testing"

	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

func TestDescContentKeyIncludesContinuation(t *testing.T) {
	label := uniqsym.New(symbol.New("accepted"))
	left := DescLabSpec{
		ConnChnlPH: symbol.New("conn"),
		ValLabQN:   label,
		ContExp:    DescCloseSpec{ContChnlPH: symbol.New("left")},
	}
	right := DescLabSpec{
		ConnChnlPH: left.ConnChnlPH,
		ValLabQN:   label,
		ContExp:    DescCloseSpec{ContChnlPH: symbol.New("right")},
	}

	leftKey, _, err := DescContentKey(left)
	if err != nil {
		t.Fatal(err)
	}
	rightKey, _, err := DescContentKey(right)
	if err != nil {
		t.Fatal(err)
	}
	if leftKey == rightKey {
		t.Fatal("different Definition continuations produced the same content key")
	}
}
