package core

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/uniqsym"
)

type CompNode struct {
	ExpCK    contkey.ADT
	SupExpCK contkey.ADT
	Kind     int16
	Spec     []byte
}

func marshalTaggedSpec(enc *jsontext.Encoder, kind string, fields map[string]any) error {
	fields["kind"] = kind
	return jsonv2.MarshalEncode(enc, fields)
}

func childCompContentKey(spec CompSpec) (int64, error) {
	key, _, err := CompContentKey(spec)
	return int64(key), err
}

func validateCompSpec(spec CompSpec) error {
	if spec == nil { return fmt.Errorf("nil pool term expression") }
	switch s := spec.(type) {
	case AcquireSpec: return validateCompSpec(s.ContExp)
	case AcceptSpec: return validateCompSpec(s.ContExp)
	case HireSpec: return validateCompSpec(s.ContExp)
	case ApplySpec: return validateCompSpec(s.ContExp)
	case ReleaseSpec, DetachSpec, FireSpec, QuitSpec, SpawnSpec:
		return nil
	case SpawnSpec2:
		return fmt.Errorf("pool SpawnSpec2 is an execution reference form and has no Definition content key")
	default: return ErrSpecTypeUnexpected(spec)
	}
}

func (s AcquireSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	return marshalTaggedSpec(enc, "acquire", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp_ck": ck})
}
func (s AcceptSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	return marshalTaggedSpec(enc, "accept", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp_ck": ck})
}
func (s HireSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	return marshalTaggedSpec(enc, "hire", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "cont_exp_ck": ck})
}
func (s ApplySpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	return marshalTaggedSpec(enc, "apply", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "cont_exp_ck": ck})
}
func (s ReleaseSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "release", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)})
}
func (s DetachSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "detach", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)})
}
func (s FireSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "fire", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN)})
}
func (s QuitSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "quit", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN)})
}
func (s SpawnSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	qns := make([]string, len(s.ProcCompQNs))
	for i, qn := range s.ProcCompQNs { qns[i] = uniqsym.ConvertToString(qn) }
	return marshalTaggedSpec(enc, "spawn", map[string]any{"proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "proc_comp_qns": qns})
}
func (s SpawnSpec2) MarshalJSONTo(enc *jsontext.Encoder) error {
	return fmt.Errorf("pool SpawnSpec2 is an execution reference form and has no Definition content key")
}

func CompContentNodes(spec CompSpec) ([]CompNode, error) {
	if err := validateCompSpec(spec); err != nil { return nil, err }
	nodes := make([]CompNode, 0)
	seen := make(map[contkey.ADT]CompNode)
	var walk func(CompSpec, contkey.ADT) error
	walk = func(current CompSpec, parent contkey.ADT) error {
		key, data, err := CompContentKey(current); if err != nil { return err }
		kind, err := compNodeKind(current); if err != nil { return err }
		node := CompNode{ExpCK: key, SupExpCK: parent, Kind: kind, Spec: data}
		if existing, ok := seen[key]; ok {
			if existing.SupExpCK != parent || existing.Kind != kind || string(existing.Spec) != string(data) {
				return fmt.Errorf("expression node %d has conflicting parent or content", key)
			}
			return nil
		}
		seen[key] = node
		nodes = append(nodes, node)
		children, err := compChildren(current); if err != nil { return err }
		for _, child := range children { if err := walk(child, key); err != nil { return err } }
		return nil
	}
	if err := walk(spec, contkey.Zero); err != nil { return nil, err }
	return nodes, nil
}

func compNodeKind(spec CompSpec) (int16, error) {
	switch spec.(type) {
	case AcquireSpec: return int16(AcquireKind), nil
	case AcceptSpec: return int16(AcceptKind), nil
	case HireSpec: return int16(HireKind), nil
	case ApplySpec: return int16(ApplyKind), nil
	case ReleaseSpec: return int16(ReleaseKind), nil
	case DetachSpec: return int16(DetachKind), nil
	case FireSpec: return int16(FireKind), nil
	case QuitSpec: return int16(QuitKind), nil
	case SpawnSpec: return int16(SpawnKind), nil
	default: return 0, ErrSpecTypeUnexpected(spec)
	}
}

func compChildren(spec CompSpec) ([]CompSpec, error) {
	switch s := spec.(type) {
	case AcquireSpec: return []CompSpec{s.ContExp}, nil
	case AcceptSpec: return []CompSpec{s.ContExp}, nil
	case HireSpec: return []CompSpec{s.ContExp}, nil
	case ApplySpec: return []CompSpec{s.ContExp}, nil
	default: return nil, nil
	}
}
