package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

func TestContentKeyDistinguishesExpressionVariants(t *testing.T) {
	acquire := termexp.AcquireSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExp:    termexp.ReleaseSpec{ConnChnlPH: symbol.New("next")},
	}
	accept := termexp.AcceptSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExp:    termexp.ReleaseSpec{ConnChnlPH: symbol.New("next")},
	}

	acquireKey, _, err := ContentKey(acquire)
	if err != nil {
		t.Fatal(err)
	}
	acceptKey, _, err := ContentKey(accept)
	if err != nil {
		t.Fatal(err)
	}
	if acquireKey == acceptKey {
		t.Fatal("different expression variants produced the same content key")
	}
}

func TestContentKeyPreservesSpawnProcessOrder(t *testing.T) {
	first := termexp.SpawnSpec{
		ProcTermQN: uniqsym.New(symbol.New("process")),
		ProcCompQNs: []uniqsym.ADT{
			uniqsym.New(symbol.New("one")),
			uniqsym.New(symbol.New("two")),
		},
	}
	second := termexp.SpawnSpec{
		ProcTermQN: uniqsym.New(symbol.New("process")),
		ProcCompQNs: []uniqsym.ADT{
			uniqsym.New(symbol.New("two")),
			uniqsym.New(symbol.New("one")),
		},
	}

	firstKey, _, err := ContentKey(first)
	if err != nil {
		t.Fatal(err)
	}
	secondKey, _, err := ContentKey(second)
	if err != nil {
		t.Fatal(err)
	}
	if firstKey == secondKey {
		t.Fatal("different spawn process order produced the same content key")
	}
}
