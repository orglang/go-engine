package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	typeexp "orglang/go-engine/internal/proc/typeexp/core"
)

func TestExpSpecFromGenerated(t *testing.T) {
	var one webapi.ProcTypeExprSpec
	if err := one.FromProcTypeOneSpec(webapi.ProcTypeOneSpec{}); err != nil {
		t.Fatal(err)
	}

	var link webapi.ProcTypeExprSpec
	if err := link.FromProcTypeLinkSpec(webapi.ProcTypeLinkSpec{
		Link: webapi.AnyTypeLinkSpec{TypeQn: "domain.type"},
	}); err != nil {
		t.Fatal(err)
	}

	var tensor webapi.ProcTypeExprSpec
	if err := tensor.FromProcTypeTensorSpec(webapi.ProcTypeTensorSpec{
		ValExpr: one,
		ContExpr: link,
	}); err != nil {
		t.Fatal(err)
	}

	var plus webapi.ProcTypeExprSpec
	if err := plus.FromProcTypePlusSpec(webapi.ProcTypePlusSpec{
		Choices: []webapi.ProcTypeChoiceSpec{{
			LabQn: "choice.a",
			ContExpr: tensor,
		}},
	}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		dto  webapi.ProcTypeExprSpec
		want any
	}{
		{name: "one", dto: one, want: typeexp.OneSpec{}},
		{name: "link", dto: link, want: typeexp.LinkSpec{}},
		{name: "tensor", dto: tensor, want: typeexp.TensorSpec{}},
		{name: "plus", dto: plus, want: typeexp.PlusSpec{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expSpecFromGenerated(tt.dto)
			if err != nil {
				t.Fatalf("convert generated expression: %v", err)
			}
			switch tt.want.(type) {
			case typeexp.OneSpec:
				if _, ok := got.(typeexp.OneSpec); !ok {
					t.Fatalf("unexpected expression type %T", got)
				}
			case typeexp.LinkSpec:
				value, ok := got.(typeexp.LinkSpec)
				if !ok || uniqsym.ConvertToString(value.TypeQN) != "domain.type" {
					t.Fatalf("unexpected link expression: %#v", got)
				}
			case typeexp.TensorSpec:
				if _, ok := got.(typeexp.TensorSpec); !ok {
					t.Fatalf("unexpected expression type %T", got)
				}
			case typeexp.PlusSpec:
				value, ok := got.(typeexp.PlusSpec)
				if !ok || len(value.Choices) != 1 {
					t.Fatalf("unexpected plus expression: %#v", got)
				}
			}
		})
	}
}

func TestExpSpecFromGeneratedRejectsUnknownKind(t *testing.T) {
	var dto webapi.ProcTypeExprSpec
	if err := dto.UnmarshalJSON([]byte(`{"kind":"unknown"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := expSpecFromGenerated(dto); err == nil {
		t.Fatal("expected unsupported expression error")
	}
}
