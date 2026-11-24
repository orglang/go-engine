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

	acquireKey, _, err := CompContentKey(acquire)
	if err != nil {
		t.Fatal(err)
	}
	acceptKey, _, err := CompContentKey(accept)
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

	firstKey, _, err := CompContentKey(first)
	if err != nil {
		t.Fatal(err)
	}
	secondKey, _, err := CompContentKey(second)
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
		spec CompSpec
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
		_, data, err := CompContentKey(tc.spec)
		if err != nil {
			t.Fatalf("%s: %v", tc.kind, err)
		}
		if !strings.Contains(string(data), `"kind":"`+tc.kind+`"`) {
			t.Fatalf("%s: missing discriminator in %s", tc.kind, data)
		}
	}
	if _, _, err := CompContentKey(SpawnSpec2{}); err == nil {
		t.Fatal("SpawnSpec2 must not be accepted as Definition content")
	}
}

func TestContentNodesUseParentKeysAndChildReferences(t *testing.T) {
	spec := AcquireSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExp: ReleaseSpec{ConnChnlPH: symbol.New("next")},
	}
	nodes, err := CompContentNodes(spec)
	if err != nil { t.Fatal(err) }
	if len(nodes) != 2 { t.Fatalf("want 2 nodes, got %d", len(nodes)) }
	if nodes[0].SupExpCK != 0 { t.Fatalf("root has parent key %d", nodes[0].SupExpCK) }
	if nodes[1].SupExpCK != nodes[0].ExpCK { t.Fatalf("child parent key = %d, want %d", nodes[1].SupExpCK, nodes[0].ExpCK) }
	if !strings.Contains(string(nodes[0].Spec), `"cont_exp_ck":`) { t.Fatalf("parent spec does not contain child key: %s", nodes[0].Spec) }
	if strings.Contains(string(nodes[0].Spec), `"cont_exp":`) { t.Fatalf("parent spec unexpectedly embeds child: %s", nodes[0].Spec) }
}

func TestContentNodesRoundTrip(t *testing.T) {
	spec := AcquireSpec{
		ConnChnlPH: symbol.New("conn"),
		ContExp: HireSpec{
			ConnChnlPH: symbol.New("worker"),
			ProcTermQN: uniqsym.New(symbol.New("designer")),
			ContExp: SpawnSpec{
				ProcTermQN: uniqsym.New(symbol.New("process")),
				ProcCompQNs: []uniqsym.ADT{uniqsym.New(symbol.New("one")), uniqsym.New(symbol.New("two"))},
			},
		},
	}
	nodes, err := CompContentNodes(spec)
	if err != nil { t.Fatal(err) }
	restored, err := BuildSpecFromNodes(nodes, nodes[0].ExpCK)
	if err != nil { t.Fatal(err) }
	wantKey, _, err := CompContentKey(spec); if err != nil { t.Fatal(err) }
	gotKey, _, err := CompContentKey(restored); if err != nil { t.Fatal(err) }
	if gotKey != wantKey { t.Fatalf("round-trip key = %d, want %d", gotKey, wantKey) }
}
