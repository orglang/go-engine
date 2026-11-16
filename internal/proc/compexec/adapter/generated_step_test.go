package adapter

import (
	"testing"

	webapi "orglang/go-engine/internal/app/web/api"
)

func TestGeneratedStepToSDKMapsProcessExpression(t *testing.T) {
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

	spec, err := generatedStepToSDK("execution-1", webapi.ProcStepSpec{ProcExpr: waitExpr})
	if err != nil {
		t.Fatalf("convert generated step: %v", err)
	}
	if spec.CompRef.CompID != "execution-1" {
		t.Errorf("unexpected computation id: got %q", spec.CompRef.CompID)
	}
	if spec.ProcExp.Wait == nil {
		t.Fatal("expected wait expression")
	}
	if spec.ProcExp.Wait.CommChnlPH != "wait-channel" {
		t.Errorf("unexpected wait channel: got %q", spec.ProcExp.Wait.CommChnlPH)
	}
	if spec.ProcExp.Wait.ContES.Close == nil || spec.ProcExp.Wait.ContES.Close.CommChnlPH != "close-channel" {
		t.Errorf("unexpected continuation: %+v", spec.ProcExp.Wait.ContES)
	}
}

func TestGeneratedStepToSDKRejectsUnsupportedLinkExpression(t *testing.T) {
	var linkExpr webapi.ProcTermExprSpec
	if err := linkExpr.FromProcTermLinkSpec(webapi.ProcTermLinkSpec{ProcQn: "p", X: "x", Ys: []string{"y"}}); err != nil {
		t.Fatalf("build generated link expression: %v", err)
	}
	if _, err := generatedStepToSDK("execution-1", webapi.ProcStepSpec{ProcExpr: linkExpr}); err == nil {
		t.Fatal("expected unsupported link expression to be rejected by SDK validation")
	}
}
