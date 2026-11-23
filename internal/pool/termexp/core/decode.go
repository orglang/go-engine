package core

import (
	jsonv2 "encoding/json/v2"
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

type storedSpec struct {
	Kind       string `json:"kind"`
	ConnChnlPH string `json:"conn_chnl_ph"`
	ProcTermQN string `json:"proc_term_qn"`
	ContExpCK  int64  `json:"cont_exp_ck"`
	ProcCompQNs []string `json:"proc_comp_qns"`
}

func BuildSpecFromNodes(nodes []ExpNode, root contkey.ADT) (ExpSpec, error) {
	byKey := make(map[contkey.ADT]ExpNode, len(nodes))
	for _, node := range nodes {
		if _, exists := byKey[node.ExpCK]; exists { return nil, fmt.Errorf("duplicate expression node %d", node.ExpCK) }
		byKey[node.ExpCK] = node
	}
	active := make(map[contkey.ADT]bool, len(nodes))
	decoded := make(map[contkey.ADT]ExpSpec, len(nodes))
	var decode func(contkey.ADT, contkey.ADT) (ExpSpec, error)
	decode = func(key, parent contkey.ADT) (ExpSpec, error) {
		node, ok := byKey[key]
		if !ok { return nil, fmt.Errorf("expression node %d does not exist", key) }
		if node.SupExpCK != parent { return nil, fmt.Errorf("expression node %d has parent %d, expected %d", key, node.SupExpCK, parent) }
		if active[key] { return nil, fmt.Errorf("expression node %d is cyclic", key) }
		if cached, ok := decoded[key]; ok { return cached, nil }
		active[key] = true
		defer delete(active, key)
		var dto storedSpec
		if err := jsonv2.Unmarshal(node.Spec, &dto); err != nil { return nil, fmt.Errorf("decode expression node %d: %w", key, err) }
		child := func() (ExpSpec, error) {
			childKey, err := contkey.ConvertFromInt(dto.ContExpCK)
			if err != nil { return nil, err }
			return decode(childKey, key)
		}
		var spec ExpSpec
		switch dto.Kind {
		case "acquire":
			ph, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			cont, err := child(); if err != nil { return nil, err }
			spec = AcquireSpec{ConnChnlPH: ph, ContExp: cont}
		case "accept":
			ph, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			cont, err := child(); if err != nil { return nil, err }
			spec = AcceptSpec{ConnChnlPH: ph, ContExp: cont}
		case "hire":
			ph, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			qn, err := uniqsym.ConvertFromString(dto.ProcTermQN); if err != nil { return nil, err }
			cont, err := child(); if err != nil { return nil, err }
			spec = HireSpec{ConnChnlPH: ph, ProcTermQN: qn, ContExp: cont}
		case "apply":
			ph, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			qn, err := uniqsym.ConvertFromString(dto.ProcTermQN); if err != nil { return nil, err }
			cont, err := child(); if err != nil { return nil, err }
			spec = ApplySpec{ConnChnlPH: ph, ProcTermQN: qn, ContExp: cont}
		case "release":
			ph, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			spec = ReleaseSpec{ConnChnlPH: ph}
		case "detach":
			ph, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			spec = DetachSpec{ConnChnlPH: ph}
		case "fire":
			ph, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			qn, err := uniqsym.ConvertFromString(dto.ProcTermQN); if err != nil { return nil, err }
			spec = FireSpec{ConnChnlPH: ph, ProcTermQN: qn}
		case "quit":
			ph, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			qn, err := uniqsym.ConvertFromString(dto.ProcTermQN); if err != nil { return nil, err }
			spec = QuitSpec{ConnChnlPH: ph, ProcTermQN: qn}
		case "spawn":
			qn, err := uniqsym.ConvertFromString(dto.ProcTermQN); if err != nil { return nil, err }
			qns, err := uniqsym.ConvertFromStrings(dto.ProcCompQNs); if err != nil { return nil, err }
			spec = SpawnSpec{ProcTermQN: qn, ProcCompQNs: qns}
		default:
			return nil, fmt.Errorf("unknown pool expression kind %q", dto.Kind)
		}
		kind, err := expNodeKind(spec); if err != nil { return nil, err }
		if kind != node.Kind { return nil, fmt.Errorf("expression node %d kind mismatch: json=%q db=%d", key, dto.Kind, node.Kind) }
		decoded[key] = spec
		return spec, nil
	}
	spec, err := decode(root, contkey.Zero)
	if err != nil { return nil, err }
	if len(decoded) != len(nodes) { return nil, fmt.Errorf("expression tree has %d unreachable nodes", len(nodes)-len(decoded)) }
	key, _, err := ContentKey(spec)
	if err != nil { return nil, err }
	if key != root { return nil, fmt.Errorf("expression root key mismatch: got %d, want %d", key, root) }
	return spec, nil
}
