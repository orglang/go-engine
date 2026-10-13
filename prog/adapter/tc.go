package adapter

import (
	"fmt"

	sdkpooltypedef "github.com/orglang/go-sdk/pool/typedef"
	sdkproctypedef "github.com/orglang/go-sdk/proc/typedef"
	sdkprog "github.com/orglang/go-sdk/prog"

	pooltypedef "orglang/go-engine/pool/typedef/adapter"
	proctypedef "orglang/go-engine/proc/typedef/adapter"
	coreprog "orglang/go-engine/prog/core"
)

func MsgToSpec(dto sdkprog.Spec) (coreprog.Spec, error) {
	spec := coreprog.Spec{}
	for _, msgExp := range dto.Exps {
		switch msgExp := msgExp.(type) {
		case sdkpooltypedef.DefSpec:
			poolDef, convErr := pooltypedef.MsgToDefSpec(msgExp)
			if convErr != nil {
				return coreprog.Spec{}, convErr
			}
			spec.Pools = append(spec.Pools, poolDef)
		case sdkproctypedef.DefSpec:
			procDef, convErr := proctypedef.MsgToDefSpec(msgExp)
			if convErr != nil {
				return coreprog.Spec{}, convErr
			}
			spec.Procs = append(spec.Procs, procDef)
		default:
			return coreprog.Spec{}, fmt.Errorf("unsupported program expression: %T", msgExp)
		}
	}
	return spec, nil
}
