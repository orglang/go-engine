package core

import (
	jsonv2 "encoding/json/v2"
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

type storedSpec struct {
	Kind         string           `json:"kind"`
	ContChnlPH   string           `json:"cont_chnl_ph"`
	ConnChnlPH   string           `json:"conn_chnl_ph"`
	ValChnlPH    string           `json:"val_chnl_ph"`
	NewChnlPH    string           `json:"new_chnl_ph"`
	ValLabQN     string           `json:"val_lab_qn"`
	ContExpCK    int64            `json:"cont_exp_ck"`
	ContExpCKs   map[string]int64 `json:"cont_exp_cks"`
	ProcTermQN   string           `json:"proc_term_qn"`
	X            string           `json:"x"`
	Ys           []string         `json:"ys"`
	ValChnlPHs   []string         `json:"val_chnl_phs"`
	NewChnlPHs   []string         `json:"new_chnl_phs"`
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
		case "close":
			ph, err := symbol.ConvertFromString(dto.ContChnlPH); if err != nil { return nil, err }
			spec = CloseSpec{ContChnlPH: ph}
		case "wait":
			ph, err := symbol.ConvertFromString(dto.ContChnlPH); if err != nil { return nil, err }
			cont, err := child(); if err != nil { return nil, err }
			spec = WaitSpec{ContChnlPH: ph, ContExp: cont}
		case "send":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			y, err := symbol.ConvertFromString(dto.ValChnlPH); if err != nil { return nil, err }
			spec = SendSpec{ConnChnlPH: x, ValChnlPH: y}
		case "recv":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			y, err := symbol.ConvertFromString(dto.NewChnlPH); if err != nil { return nil, err }
			cont, err := child(); if err != nil { return nil, err }
			spec = RecvSpec{ConnChnlPH: x, NewChnlPH: y, ContExp: cont}
		case "lab":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			label, err := uniqsym.ConvertFromString(dto.ValLabQN); if err != nil { return nil, err }
			cont, err := child(); if err != nil { return nil, err }
			spec = LabSpec{ConnChnlPH: x, ValLabQN: label, ContExp: cont}
		case "case":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			branches := make(map[uniqsym.ADT]ExpSpec, len(dto.ContExpCKs))
			for label, rawCK := range dto.ContExpCKs {
				qn, err := uniqsym.ConvertFromString(label); if err != nil { return nil, err }
				childKey, err := contkey.ConvertFromInt(rawCK); if err != nil { return nil, err }
				cont, err := decode(childKey, key); if err != nil { return nil, err }
				branches[qn] = cont
			}
			spec = CaseSpec{ConnChnlPH: x, ContExps: branches}
		case "link":
			qn, err := uniqsym.ConvertFromString(dto.ProcTermQN); if err != nil { return nil, err }
			x, err := identity.ConvertFromString(dto.X); if err != nil { return nil, err }
			ys := make([]identity.ADT, len(dto.Ys))
			for i, raw := range dto.Ys { ys[i], err = identity.ConvertFromString(raw); if err != nil { return nil, err } }
			spec = LinkSpec{ProcTermQN: qn, X: x, Ys: ys}
		case "fwd":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			y, err := symbol.ConvertFromString(dto.ContChnlPH); if err != nil { return nil, err }
			spec = FwdSpec{ConnChnlPH: x, ContChnlPH: y}
		case "call":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			y, err := symbol.ConvertFromString(dto.NewChnlPH); if err != nil { return nil, err }
			qn, err := uniqsym.ConvertFromString(dto.ProcTermQN); if err != nil { return nil, err }
			chs := make([]symbol.ADT, len(dto.ValChnlPHs))
			for i, raw := range dto.ValChnlPHs { chs[i], err = symbol.ConvertFromString(raw); if err != nil { return nil, err } }
			cont, err := child(); if err != nil { return nil, err }
			spec = CallSpec{ConnChnlPH: x, NewChnlPH: y, ProcTermQN: qn, ValChnlPHs: chs, ContExp: cont}
		case "spawn":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			qn, err := uniqsym.ConvertFromString(dto.ProcTermQN); if err != nil { return nil, err }
			chs := make([]symbol.ADT, len(dto.NewChnlPHs))
			for i, raw := range dto.NewChnlPHs { chs[i], err = symbol.ConvertFromString(raw); if err != nil { return nil, err } }
			cont, err := child(); if err != nil { return nil, err }
			spec = SpawnSpec{ConnChnlPH: x, ProcTermQN: qn, NewChnlPHs: chs, ContExp: cont}
		case "acquire":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			cont, err := child(); if err != nil { return nil, err }
			spec = AcqureSpec{ConnChnlPH: x, ContExp: cont}
		case "accept":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			cont, err := child(); if err != nil { return nil, err }
			spec = AcceptSpec{ConnChnlPH: x, ContExp: cont}
		case "detach":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			spec = DetachSpec{ConnChnlPH: x}
		case "release":
			x, err := symbol.ConvertFromString(dto.ConnChnlPH); if err != nil { return nil, err }
			spec = ReleaseSpec{ConnChnlPH: x}
		default:
			return nil, fmt.Errorf("unknown process expression kind %q", dto.Kind)
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
