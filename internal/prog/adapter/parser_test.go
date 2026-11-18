package adapter

import (
	"testing"

	pooltypeexp "orglang/go-engine/internal/pool/typeexp/core"
	proctypexp "orglang/go-engine/internal/proc/typeexp/core"
)

func TestParseProgramAndConvertToCore(t *testing.T) {
	const source = `pool type PoolA { plus (ProcA ProcB) { one } }
proc type ProcA { tensor { one } { plus { @ok { one } @fail { link PoolA } } } }`

	ast, err := parseProgram(source)
	if err != nil {
		t.Fatalf("parse program: %v", err)
	}
	spec, err := programASTToCore(ast)
	if err != nil {
		t.Fatalf("convert program: %v", err)
	}
	if len(spec.Pools) != 1 || len(spec.Procs) != 1 {
		t.Fatalf("unexpected definition counts: pools=%d procs=%d", len(spec.Pools), len(spec.Procs))
	}
	if _, ok := spec.Pools[0].TypeExp.(pooltypeexp.PlusSpec); !ok {
		t.Fatalf("unexpected pool expression type: %T", spec.Pools[0].TypeExp)
	}
	procExp, ok := spec.Procs[0].TypeExp.(proctypexp.TensorSpec)
	if !ok {
		t.Fatalf("unexpected process expression type: %T", spec.Procs[0].TypeExp)
	}
	if _, ok := procExp.Cont.(proctypexp.PlusSpec); !ok {
		t.Fatalf("unexpected process continuation type: %T", procExp.Cont)
	}
}

func TestParseProgramRejectsUnknownExpression(t *testing.T) {
	if _, err := parseProgram("unknown type Foo { one }"); err == nil {
		t.Fatal("expected parser error for unknown declaration")
	}
}

func TestProgramConversionRejectsDuplicateProcessLabels(t *testing.T) {
	ast, err := parseProgram("proc type Foo { plus { @same { one } @same { one } } }")
	if err != nil {
		t.Fatalf("parse program: %v", err)
	}
	if _, err := programASTToCore(ast); err == nil {
		t.Fatal("expected duplicate-label conversion error")
	}
}
