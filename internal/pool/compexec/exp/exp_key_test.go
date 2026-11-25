package exp

import (
	"strings"
	"testing"

	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

func TestExpContentKeyDistinguishesExpressionVariants(t *testing.T) {
	acquire := AcquireSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExp:    ReleaseSpec{ConnChnlPH: symbol.New("next")},
	}
	accept := AcceptSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExp:    ReleaseSpec{ConnChnlPH: symbol.New("next")},
	}

	acquireKey, _, err := ExpContentKey(acquire)
	if err != nil {
		t.Fatal(err)
	}
	acceptKey, _, err := ExpContentKey(accept)
	if err != nil {
		t.Fatal(err)
	}
	if acquireKey == acceptKey {
		t.Fatal("different expression variants produced the same content key")
	}
}

func TestExpContentKeyPreservesSpawnProcessOrder(t *testing.T) {
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

	firstKey, _, err := ExpContentKey(first)
	if err != nil {
		t.Fatal(err)
	}
	secondKey, _, err := ExpContentKey(second)
	if err != nil {
		t.Fatal(err)
	}
	if firstKey == secondKey {
		t.Fatal("different spawn process order produced the same content key")
	}
}

func TestExpContentKeySerializesEveryVariant(t *testing.T) {
	cont := ReleaseSpec{ConnChnlPH: symbol.New("next")}
	procTermQN := uniqsym.New(symbol.New("process"))
	specs := []struct {
		spec ExpSpec
		kind string
	}{
		{AcquireSpec{ConnChnlPH: symbol.New("conn"), ContExp: cont}, "acquire"},
		{AcceptSpec{ConnChnlPH: symbol.New("conn"), ContExp: cont}, "accept"},
		{HireSpec{ConnChnlPH: symbol.New("conn"), ProcTermQN: procTermQN, ContExp: cont}, "hire"},
		{ApplySpec{ConnChnlPH: symbol.New("conn"), ProcTermQN: procTermQN, ContExp: cont}, "apply"},
		{ReleaseSpec{ConnChnlPH: symbol.New("conn")}, "release"},
		{DetachSpec{ConnChnlPH: symbol.New("conn")}, "detach"},
		{FireSpec{ConnChnlPH: symbol.New("conn"), ProcTermQN: procTermQN}, "fire"},
		{QuitSpec{ConnChnlPH: symbol.New("conn"), ProcTermQN: procTermQN}, "quit"},
		{SpawnSpec{ProcTermQN: procTermQN, ProcCompQNs: []uniqsym.ADT{procTermQN}}, "spawn"},
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
	if _, _, err := ExpContentKey(SpawnSpec2{}); err == nil {
		t.Fatal("SpawnSpec2 must not be accepted as Definition content")
	}
}
