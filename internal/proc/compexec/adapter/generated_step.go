package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"

	sdkcompsem "github.com/orglang/go-sdk/adt/compsem"
	sdkcompstep "github.com/orglang/go-sdk/proc/compstep"
	"github.com/labstack/echo/v4"

	webapi "orglang/go-engine/internal/app/web/api"
	compstep "orglang/go-engine/internal/proc/compstep/adapter"
)

func generatedStepToSDK(id string, dto webapi.ProcStepSpec) (sdkcompstep.StepSpec, error) {
	var spec sdkcompstep.StepSpec
	spec.CompRef = sdkcompsem.SemRef{CompID: id}

	raw, err := json.Marshal(dto.ProcExpr)
	if err != nil {
		return spec, fmt.Errorf("marshal generated process expression: %w", err)
	}
	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return spec, fmt.Errorf("decode generated process expression: %w", err)
	}
	renameProcessExpressionFields(value)
	raw, err = json.Marshal(value)
	if err != nil {
		return spec, fmt.Errorf("marshal normalized process expression: %w", err)
	}
	if err := json.Unmarshal(raw, &spec.ProcExp); err != nil {
		return spec, fmt.Errorf("decode SDK process expression: %w", err)
	}
	if err := spec.ProcExp.Validate(); err != nil {
		return spec, fmt.Errorf("validate process expression: %w", err)
	}
	return spec, nil
}

func renameProcessExpressionFields(value interface{}) {
	switch node := value.(type) {
	case map[string]interface{}:
		for key, child := range node {
			renameProcessExpressionFields(child)
			switch key {
			case "conn_ph":
				node["comm_ph"] = child
				delete(node, key)
			case "cont_expr":
				node["cont_es"] = child
				delete(node, key)
			}
		}
	case []interface{}:
		for _, child := range node {
			renameProcessExpressionFields(child)
		}
	}
}

func (h *controllerEcho) StepProcessExecution(c echo.Context, id string) error {
	var dto webapi.StepProcessExecutionJSONRequestBody
	if err := c.Bind(&dto); err != nil {
		h.log.Error("binding generated process step failed", "error", err)
		return err
	}
	spec, err := generatedStepToSDK(id, dto)
	if err != nil {
		h.log.Error("converting generated process step failed", "error", err)
		return err
	}
	coreSpec, err := compstep.MsgToStepSpec(spec)
	if err != nil {
		h.log.Error("converting process step to core model failed", "error", err)
		return err
	}
	if err := h.api.Take(coreSpec); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

