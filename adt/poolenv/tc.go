package poolenv

import (
	"orglang/go-engine/adt/poolexp"
	"orglang/go-engine/adt/uniqsym"
)

func ConvertExpToEnv(s poolexp.ExpSpec) EnvSpec {
	procDescQNs := make([]uniqsym.ADT, 0, 10)
	procImplQNs := make([]uniqsym.ADT, 0, 10)
	switch spec := s.(type) {
	case poolexp.HireSpec:
		procDescQNs = append(procDescQNs, spec.ProcDescQN)
	case poolexp.ApplySpec:
		procDescQNs = append(procDescQNs, spec.ProcDescQN)
	case poolexp.FireSpec:
		procDescQNs = append(procDescQNs, spec.ProcDescQN)
	case poolexp.QuitSpec:
		procDescQNs = append(procDescQNs, spec.ProcDescQN)
	case poolexp.SpawnSpec:
		procDescQNs = append(procDescQNs, spec.ProcDescQN)
		procImplQNs = append(procImplQNs, spec.ProcImplQNs...)
	default:
		panic(poolexp.ErrSpecTypeUnexpected(s))
	}
	return EnvSpec{ProcDescQNs: procDescQNs, ProcImplQNs: procImplQNs}
}
