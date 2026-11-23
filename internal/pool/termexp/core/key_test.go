package core

import (
	"testing"
	"strings"

	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

func TestContentKeyDistinguishesExpressionVariants(t *testing.T) {
	acquire := AcquireSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExp:    ReleaseSpec{ConnChnlPH: symbol.New("next")},
	}
	accept := AcceptSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExp:    ReleaseSpec{ConnChnlPH: symbol.New("next")},
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
	first := SpawnSpec{
		ProcTermQN: uniqsym.New(symbol.New("process")),
		ProcCompQNs: []uniqsym.ADT{
			uniqsym.New(symbol.New("one")),
			uniqsym.New(symbol.New("two")),
		},
	}
	second := SpawnSpec{
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

func TestContentKeySerializesEveryDefinitionVariant(t *testing.T) {
	cont := ReleaseSpec{}
	specs := []struct {
		spec ExpSpec
		kind string
	}{
		{AcquireSpec{ContExp: cont}, "acquire"},
		{AcceptSpec{ContExp: cont}, "accept"},
		{HireSpec{ContExp: cont}, "hire"},
		{ApplySpec{ContExp: cont}, "apply"},
		{ReleaseSpec{}, "release"},
		{DetachSpec{}, "detach"},
		{FireSpec{}, "fire"},
		{QuitSpec{}, "quit"},
		{SpawnSpec{}, "spawn"},
	}
	for _, tc := range specs {
		_, data, err := ContentKey(tc.spec)
		if err != nil {
			t.Fatalf("%s: %v", tc.kind, err)
		}
		if !strings.Contains(string(data), `"kind":"`+tc.kind+`"`) {
			t.Fatalf("%s: missing discriminator in %s", tc.kind, data)
		}
	}
	if _, _, err := ContentKey(SpawnSpec2{}); err == nil {
		t.Fatal("SpawnSpec2 must not be accepted as Definition content")
	}
}
