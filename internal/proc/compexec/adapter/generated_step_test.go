package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	webapi "orglang/go-engine/internal/app/web/api"
	termexpcore "orglang/go-engine/internal/proc/termexp/core"
)

const testExecutionID = "9m4e2mr0ui3e8a215n4g"

func TestGeneratedStepToCoreMapsNestedProcessExpression(t *testing.T) {
	var closeExpr webapi.ProcTermExprSpec
	if err := closeExpr.FromProcTermCloseSpec(webapi.ProcTermCloseSpec{ConnPh: "close-channel"}); err != nil {
		t.Fatalf("build generated close expression: %v", err)
	}
	var waitExpr webapi.ProcTermExprSpec
	if err := waitExpr.FromProcTermWaitSpec(webapi.ProcTermWaitSpec{
		ConnPh:   "wait-channel",
		ContExpr: closeExpr,
	}); err != nil {
		t.Fatalf("build generated wait expression: %v", err)
	}

	spec, err := generatedStepToCore(testExecutionID, webapi.ProcStepSpec{ProcExpr: waitExpr})
	if err != nil {
		t.Fatalf("convert generated step: %v", err)
	}
	if got := identity.ConvertToString(spec.CompRef.CompID); got != testExecutionID {
		t.Errorf("unexpected computation id: got %q", got)
	}
	wait, ok := spec.ProcExp.(termexpcore.WaitSpec)
	if !ok {
		t.Fatalf("expected core wait expression, got %T", spec.ProcExp)
	}
	if got := symbol.ConvertToString(wait.ContChnlPH); got != "wait-channel" {
		t.Errorf("unexpected wait channel: got %q", got)
	}
	close, ok := wait.ContExp.(termexpcore.CloseSpec)
	if !ok {
		t.Fatalf("expected core close continuation, got %T", wait.ContExp)
	}
	if got := symbol.ConvertToString(close.ContChnlPH); got != "close-channel" {
		t.Errorf("unexpected close channel: got %q", got)
	}
}

func TestGeneratedStepToCoreMapsLinkExpression(t *testing.T) {
	var linkExpr webapi.ProcTermExprSpec
	if err := linkExpr.FromProcTermLinkSpec(webapi.ProcTermLinkSpec{
		ProcQn: "process",
		X:      testExecutionID,
		Ys:     []string{testExecutionID},
	}); err != nil {
		t.Fatalf("build generated link expression: %v", err)
	}

	spec, err := generatedStepToCore(testExecutionID, webapi.ProcStepSpec{ProcExpr: linkExpr})
	if err != nil {
		t.Fatalf("convert generated link expression: %v", err)
	}
	link, ok := spec.ProcExp.(termexpcore.LinkSpec)
	if !ok {
		t.Fatalf("expected core link expression, got %T", spec.ProcExp)
	}
	if got := identity.ConvertToString(link.X); got != testExecutionID {
		t.Errorf("unexpected linked identity: got %q", got)
	}
	if len(link.Ys) != 1 || identity.ConvertToString(link.Ys[0]) != testExecutionID {
		t.Errorf("unexpected linked identities: %+v", link.Ys)
	}
}

func TestGeneratedStepToCoreRejectsInvalidExecutionID(t *testing.T) {
	var closeExpr webapi.ProcTermExprSpec
	if err := closeExpr.FromProcTermCloseSpec(webapi.ProcTermCloseSpec{ConnPh: "close-channel"}); err != nil {
		t.Fatalf("build generated close expression: %v", err)
	}
	if _, err := generatedStepToCore("not-an-xid", webapi.ProcStepSpec{ProcExpr: closeExpr}); err == nil {
		t.Fatal("expected invalid execution id to be rejected")
	}
}
