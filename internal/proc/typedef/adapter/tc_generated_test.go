package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

func TestExpSpecFromGenerated(t *testing.T) {
	var one webapi.ProcTypeExprSpec
	err := one.FromProcTypeOneSpec(webapi.ProcTypeOneSpec{})
	if err != nil {
		t.Fatal(err)
	}

	linkSpec := webapi.ProcTypeLinkSpec{
		Link: webapi.AnyTypeLinkSpec{TypeQn: "domain.type"},
	}
	var link webapi.ProcTypeExprSpec
	err = link.FromProcTypeLinkSpec(linkSpec)
	if err != nil {
		t.Fatal(err)
	}

	tensorSpec := webapi.ProcTypeTensorSpec{
		ValExpr:  one,
		ContExpr: link,
	}
	var tensor webapi.ProcTypeExprSpec
	err = tensor.FromProcTypeTensorSpec(tensorSpec)
	if err != nil {
		t.Fatal(err)
	}

	plusSpec := webapi.ProcTypePlusSpec{
		Choices: []webapi.ProcTypeChoiceSpec{{
			LabQn:    "choice.a",
			ContExpr: tensor,
		}},
	}
	var plus webapi.ProcTypeExprSpec
	err = plus.FromProcTypePlusSpec(plusSpec)
	if err != nil {
		t.Fatal(err)
	}

	var up webapi.ProcTypeExprSpec
	if err = up.FromProcTypeUpSpec(webapi.ProcTypeUpSpec{ContExpr: link}); err != nil {
		t.Fatal(err)
	}
	var down webapi.ProcTypeExprSpec
	if err = down.FromProcTypeDownSpec(webapi.ProcTypeDownSpec{ContExpr: one}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		dto  webapi.ProcTypeExprSpec
		want string
	}{
		{name: "one", dto: one, want: "one"},
		{name: "link", dto: link, want: "link"},
		{name: "tensor", dto: tensor, want: "tensor"},
		{name: "plus", dto: plus, want: "plus"},
		{name: "up", dto: up, want: "up"},
		{name: "down", dto: down, want: "down"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exprSpecFromGenerated(tt.dto)
			if err != nil {
				t.Fatalf("convert generated expression: %v", err)
			}
			switch tt.want {
			case "one":
				if _, ok := got.(typedef.OneSpec); !ok {
					t.Fatalf("unexpected expression type %T", got)
				}
			case "link":
				value, ok := got.(typedef.LinkSpec)
				if !ok || uniqsym.ConvertToString(value.TypeQN) != "domain.type" {
					t.Fatalf("unexpected link expression: %#v", got)
				}
			case "tensor":
				if _, ok := got.(typedef.TensorSpec); !ok {
					t.Fatalf("unexpected expression type %T", got)
				}
			case "plus":
				value, ok := got.(typedef.PlusSpec)
				if !ok || len(value.Choices) != 1 {
					t.Fatalf("unexpected plus expression: %#v", got)
				}
			case "up":
				if _, ok := got.(typedef.UpSpec); !ok {
					t.Fatalf("unexpected expression type %T", got)
				}
			case "down":
				if _, ok := got.(typedef.DownSpec); !ok {
					t.Fatalf("unexpected expression type %T", got)
				}
			default:
				t.Fatalf("unknown expected expression type %q", tt.want)
			}
		})
	}
}

func TestExpSpecFromGeneratedRejectsUnknownKind(t *testing.T) {
	var dto webapi.ProcTypeExprSpec
	if err := dto.UnmarshalJSON([]byte(`{"kind":"unknown"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := exprSpecFromGenerated(dto); err == nil {
		t.Fatal("expected unsupported expression error")
	}
}
