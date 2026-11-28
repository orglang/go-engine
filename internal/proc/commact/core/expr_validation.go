package core

import (
	"errors"
	"fmt"

	"orglang/go-engine/internal/adt/uniqsym"
)

func validateExprSpec(spec ExprSpec) error {
	if spec == nil {
		return errors.New("nil process term expression")
	}
	switch s := spec.(type) {
	case CloseSpec, SendSpec, LinkSpec, FwdSpec, DetachSpec, ReleaseSpec:
		return nil
	case WaitSpec:
		return validateExprSpec(s.ContExpr)
	case RecvSpec:
		return validateExprSpec(s.ContExpr)
	case LabSpec:
		return validateExprSpec(s.ContExpr)
	case CaseSpec:
		for label, child := range s.ContExprs {
			if err := validateExprSpec(child); err != nil {
				return fmt.Errorf("case[%s]: %w", uniqsym.ConvertToString(label), err)
			}
		}
		return nil
	case CallSpec:
		return validateExprSpec(s.ContExpr)
	case SpawnSpec:
		return validateExprSpec(s.ContExpr)
	case AcqureSpec:
		return validateExprSpec(s.ContExpr)
	case AcceptSpec:
		return validateExprSpec(s.ContExpr)
	default:
		return ErrExprTypeUnexpected(spec)
	}
}
