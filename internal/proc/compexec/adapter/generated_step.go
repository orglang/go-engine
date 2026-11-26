package adapter

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	webapi "orglang/go-engine/internal/app/web/api"
	compexp "orglang/go-engine/internal/proc/compexec/core"
	compstepcore "orglang/go-engine/internal/proc/compstep/core"
)

func generatedStepToCore(id string, dto webapi.ProcStepSpec) (compstepcore.StepSpec, error) {
	compID, err := identity.ConvertFromString(id)
	if err != nil {
		return compstepcore.StepSpec{}, fmt.Errorf("convert process execution id: %w", err)
	}
	expr, err := generatedProcessExpressionToCore(dto.ProcExpr)
	if err != nil {
		return compstepcore.StepSpec{}, err
	}
	return compstepcore.StepSpec{
		CompRef: comproot.SemRef{CompID: compID},
		ProcExp: expr,
	}, nil
}

func generatedProcessExpressionToCore(dto webapi.ProcTermExprSpec) (compexp.ExpSpec, error) {
	kind, err := dto.Discriminator()
	if err != nil {
		return nil, fmt.Errorf("read generated process expression kind: %w", err)
	}

	switch kind {
	case "close":
		spec, err := dto.AsProcTermCloseSpec()
		if err != nil {
			return nil, err
		}
		channel, err := symbol.ConvertFromString(spec.ConnPh)
		if err != nil {
			return nil, err
		}
		return compexp.CloseSpec{ContChnlPH: channel}, nil
	case "wait":
		spec, err := dto.AsProcTermWaitSpec()
		if err != nil {
			return nil, err
		}
		channel, err := symbol.ConvertFromString(spec.ConnPh)
		if err != nil {
			return nil, err
		}
		cont, err := generatedProcessExpressionToCore(spec.ContExpr)
		if err != nil {
			return nil, err
		}
		return compexp.WaitSpec{ContChnlPH: channel, ContExp: cont}, nil
	case "send":
		spec, err := dto.AsProcTermSendSpec()
		if err != nil {
			return nil, err
		}
		conn, err := symbol.ConvertFromString(spec.ConnPh)
		if err != nil {
			return nil, err
		}
		value, err := symbol.ConvertFromString(spec.ValPh)
		if err != nil {
			return nil, err
		}
		return compexp.SendSpec{ConnChnlPH: conn, ValChnlPH: value}, nil
	case "recv":
		spec, err := dto.AsProcTermRecvSpec()
		if err != nil {
			return nil, err
		}
		conn, err := symbol.ConvertFromString(spec.ConnPh)
		if err != nil {
			return nil, err
		}
		bind, err := symbol.ConvertFromString(spec.BindPh)
		if err != nil {
			return nil, err
		}
		cont, err := generatedProcessExpressionToCore(spec.ContExpr)
		if err != nil {
			return nil, err
		}
		return compexp.RecvSpec{ConnChnlPH: conn, NewChnlPH: bind, ContExp: cont}, nil
	case "lab":
		spec, err := dto.AsProcTermLabSpec()
		if err != nil {
			return nil, err
		}
		conn, err := symbol.ConvertFromString(spec.ConnPh)
		if err != nil {
			return nil, err
		}
		label, err := uniqsym.ConvertFromString(spec.PatternQn)
		if err != nil {
			return nil, err
		}
		return compexp.LabSpec{ConnChnlPH: conn, ValLabQN: label}, nil
	case "case":
		spec, err := dto.AsProcTermCaseSpec()
		if err != nil {
			return nil, err
		}
		conn, err := symbol.ConvertFromString(spec.ConnPh)
		if err != nil {
			return nil, err
		}
		branches := make(map[uniqsym.ADT]compexp.ExpSpec, len(spec.ContBses))
		for _, branch := range spec.ContBses {
			label, err := uniqsym.ConvertFromString(branch.PatternQn)
			if err != nil {
				return nil, err
			}
			cont, err := generatedProcessExpressionToCore(branch.ContExpr)
			if err != nil {
				return nil, err
			}
			branches[label] = cont
		}
		return compexp.CaseSpec{ConnChnlPH: conn, ContExps: branches}, nil
	case "link":
		spec, err := dto.AsProcTermLinkSpec()
		if err != nil {
			return nil, err
		}
		proc, err := uniqsym.ConvertFromString(spec.ProcQn)
		if err != nil {
			return nil, err
		}
		x, err := identity.ConvertFromString(spec.X)
		if err != nil {
			return nil, err
		}
		ys := make([]identity.ADT, len(spec.Ys))
		for i, value := range spec.Ys {
			ys[i], err = identity.ConvertFromString(value)
			if err != nil {
				return nil, err
			}
		}
		return compexp.LinkSpec{ProcTermQN: proc, X: x, Ys: ys}, nil
	case "fwd":
		spec, err := dto.AsProcTermFwdSpec()
		if err != nil {
			return nil, err
		}
		conn, err := symbol.ConvertFromString(spec.ConnPh)
		if err != nil {
			return nil, err
		}
		cont, err := symbol.ConvertFromString(spec.ContPh)
		if err != nil {
			return nil, err
		}
		return compexp.FwdSpec{ConnChnlPH: conn, ContChnlPH: cont}, nil
	case "call":
		spec, err := dto.AsProcTermCallSpec()
		if err != nil {
			return nil, err
		}
		bind, err := symbol.ConvertFromString(spec.BindPh)
		if err != nil {
			return nil, err
		}
		proc, err := uniqsym.ConvertFromString(spec.ProcQn)
		if err != nil {
			return nil, err
		}
		values := make([]symbol.ADT, len(spec.ValPhs))
		for i, value := range spec.ValPhs {
			values[i], err = symbol.ConvertFromString(value)
			if err != nil {
				return nil, err
			}
		}
		cont, err := generatedProcessExpressionToCore(spec.ContExpr)
		if err != nil {
			return nil, err
		}
		return compexp.CallSpec{NewChnlPH: bind, ProcTermQN: proc, ValChnlPHs: values, ContExp: cont}, nil
	default:
		return nil, fmt.Errorf("unsupported process expression kind %q", kind)
	}
}

func (h *controllerEcho) StepProcessExecution(c *echo.Context, id string) error {
	var dto webapi.StepProcessExecutionJSONRequestBody
	if err := c.Bind(&dto); err != nil {
		h.log.Error("binding generated process step failed", "error", err)
		return err
	}
	spec, err := generatedStepToCore(id, dto)
	if err != nil {
		h.log.Error("converting generated process step failed", "error", err)
		return err
	}
	if err := h.api.Take(c.Request().Context(), spec); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
