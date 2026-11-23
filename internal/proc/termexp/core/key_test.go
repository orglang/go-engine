package core

import (
	"testing"
	"strings"

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

func TestContentKeySerializesEveryDefinitionVariant(t *testing.T) {
	cont := CloseSpec{ContChnlPH: symbol.New("next")}
	specs := []struct {
		spec ExpSpec
		kind string
	}{
		{CloseSpec{}, "close"},
		{WaitSpec{ContExp: cont}, "wait"},
		{SendSpec{}, "send"},
		{RecvSpec{ContExp: cont}, "recv"},
		{LabSpec{ContExp: cont}, "lab"},
		{CaseSpec{ContExps: map[uniqsym.ADT]ExpSpec{}}, "case"},
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
		_, data, err := ContentKey(tc.spec)
		if err != nil {
			t.Fatalf("%s: %v", tc.kind, err)
		}
		if !strings.Contains(string(data), `"kind":"`+tc.kind+`"`) {
			t.Fatalf("%s: missing discriminator in %s", tc.kind, data)
		}
	}
}

func TestContentNodesUseParentKeysAndChildReferences(t *testing.T) {
	spec := LabSpec{
		ConnChnlPH: symbol.New("conn"),
		ValLabQN: uniqsym.New(symbol.New("accepted")),
		ContExp: CloseSpec{ContChnlPH: symbol.New("next")},
	}
	nodes, err := ContentNodes(spec)
	if err != nil { t.Fatal(err) }
	if len(nodes) != 2 { t.Fatalf("want 2 nodes, got %d", len(nodes)) }
	if nodes[0].SupExpCK != 0 { t.Fatalf("root has parent key %d", nodes[0].SupExpCK) }
	if nodes[1].SupExpCK != nodes[0].ExpCK { t.Fatalf("child parent key = %d, want %d", nodes[1].SupExpCK, nodes[0].ExpCK) }
	if !strings.Contains(string(nodes[0].Spec), `"cont_exp_ck":`) { t.Fatalf("parent spec does not contain child key: %s", nodes[0].Spec) }
	if strings.Contains(string(nodes[0].Spec), `"cont_exp":`) { t.Fatalf("parent spec unexpectedly embeds child: %s", nodes[0].Spec) }
}
