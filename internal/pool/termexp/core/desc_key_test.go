package core

import (
	"testing"

	"orglang/go-engine/internal/adt/symbol"
)

func TestDescContentKeyDistinguishesDefinitionVariants(t *testing.T) {
	continuation := DescReleaseSpec{ConnChnlPH: symbol.New("next")}
	acquire := DescAcquireSpec{ConnChnlPH: symbol.New("conn"), ContExp: continuation}
	accept := DescAcceptSpec{ConnChnlPH: symbol.New("conn"), ContExp: continuation}

	acquireKey, _, err := DescContentKey(acquire)
	if err != nil {
		t.Fatal(err)
	}
	acceptKey, _, err := DescContentKey(accept)
	if err != nil {
		t.Fatal(err)
	}
	if acquireKey == acceptKey {
		t.Fatal("different Definition variants produced the same content key")
	}
}
