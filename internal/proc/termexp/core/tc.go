package core

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"fmt"
	"sort"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
)

const (
	CloseKind int16 = iota + 1
	WaitKind
	SendKind
	RecvKind
	LabKind
	CaseKind
	LinkKind
	FwdKind
	CallKind
	SpawnKind
	AcquireKind
	AcceptKind
	DetachKind
	ReleaseKind
)

// CompNode is one immutable node in a persisted Definition tree.
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
	if spec == nil { return fmt.Errorf("nil process term expression") }
	switch s := spec.(type) {
	case CloseSpec, SendSpec, LinkSpec, FwdSpec, DetachSpec, ReleaseSpec:
		return nil
	case WaitSpec: return validateCompSpec(s.ContExp)
	case RecvSpec: return validateCompSpec(s.ContExp)
	case LabSpec: return validateCompSpec(s.ContExp)
	case CaseSpec:
		for label, child := range s.ContExps {
			if err := validateCompSpec(child); err != nil { return fmt.Errorf("case[%s]: %w", uniqsym.ConvertToString(label), err) }
		}
		return nil
	case CallSpec: return validateCompSpec(s.ContExp)
	case SpawnSpec: return validateCompSpec(s.ContExp)
	case AcqureSpec: return validateCompSpec(s.ContExp)
	case AcceptSpec: return validateCompSpec(s.ContExp)
	default: return ErrExpTypeUnexpected(spec)
	}
}

func (s CloseSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "close", map[string]any{"cont_chnl_ph": string(s.ContChnlPH)})
}
func (s WaitSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	return marshalTaggedSpec(enc, "wait", map[string]any{"cont_chnl_ph": string(s.ContChnlPH), "cont_exp_ck": ck})
}
func (s SendSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "send", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "val_chnl_ph": string(s.ValChnlPH)})
}
func (s RecvSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	return marshalTaggedSpec(enc, "recv", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "new_chnl_ph": string(s.NewChnlPH), "cont_exp_ck": ck})
}
func (s LabSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	return marshalTaggedSpec(enc, "lab", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "val_lab_qn": uniqsym.ConvertToString(s.ValLabQN), "cont_exp_ck": ck})
}
func (s CaseSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	branches := make(map[string]int64, len(s.ContExps))
	for qn, child := range s.ContExps {
		name := uniqsym.ConvertToString(qn)
		if _, exists := branches[name]; exists { return fmt.Errorf("duplicate case label %q", name) }
		ck, err := childCompContentKey(child); if err != nil { return fmt.Errorf("case[%s]: %w", name, err) }
		branches[name] = ck
	}
	return marshalTaggedSpec(enc, "case", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp_cks": branches})
}
func (s LinkSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ys := make([]string, len(s.Ys))
	for i, id := range s.Ys { ys[i] = identity.ConvertToString(id) }
	return marshalTaggedSpec(enc, "link", map[string]any{"proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "x": identity.ConvertToString(s.X), "ys": ys})
}
func (s FwdSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "fwd", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_chnl_ph": string(s.ContChnlPH)})
}
func (s CallSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	chs := make([]string, len(s.ValChnlPHs))
	for i, ch := range s.ValChnlPHs { chs[i] = string(ch) }
	return marshalTaggedSpec(enc, "call", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "new_chnl_ph": string(s.NewChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "val_chnl_phs": chs, "cont_exp_ck": ck})
}
func (s SpawnSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	chs := make([]string, len(s.NewChnlPHs))
	for i, ch := range s.NewChnlPHs { chs[i] = string(ch) }
	return marshalTaggedSpec(enc, "spawn", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "new_chnl_phs": chs, "cont_exp_ck": ck})
}
func (s AcqureSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	return marshalTaggedSpec(enc, "acquire", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp_ck": ck})
}
func (s AcceptSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	ck, err := childCompContentKey(s.ContExp); if err != nil { return err }
	return marshalTaggedSpec(enc, "accept", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp_ck": ck})
}
func (s DetachSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "detach", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)})
}
func (s ReleaseSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "release", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)})
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
	case CloseSpec: return CloseKind, nil
	case WaitSpec: return WaitKind, nil
	case SendSpec: return SendKind, nil
	case RecvSpec: return RecvKind, nil
	case LabSpec: return LabKind, nil
	case CaseSpec: return CaseKind, nil
	case LinkSpec: return LinkKind, nil
	case FwdSpec: return FwdKind, nil
	case CallSpec: return CallKind, nil
	case SpawnSpec: return SpawnKind, nil
	case AcqureSpec: return AcquireKind, nil
	case AcceptSpec: return AcceptKind, nil
	case DetachSpec: return DetachKind, nil
	case ReleaseSpec: return ReleaseKind, nil
	default: return 0, ErrExpTypeUnexpected(spec)
	}
}

func compChildren(spec CompSpec) ([]CompSpec, error) {
	switch s := spec.(type) {
	case WaitSpec: return []CompSpec{s.ContExp}, nil
	case RecvSpec: return []CompSpec{s.ContExp}, nil
	case LabSpec: return []CompSpec{s.ContExp}, nil
	case CaseSpec:
		labels := make([]uniqsym.ADT, 0, len(s.ContExps))
		for label := range s.ContExps { labels = append(labels, label) }
		sort.Slice(labels, func(i,j int) bool { return uniqsym.ConvertToString(labels[i]) < uniqsym.ConvertToString(labels[j]) })
		children := make([]CompSpec, 0, len(labels))
		for _, label := range labels { children = append(children, s.ContExps[label]) }
		return children, nil
	case CallSpec: return []CompSpec{s.ContExp}, nil
	case SpawnSpec: return []CompSpec{s.ContExp}, nil
	case AcqureSpec: return []CompSpec{s.ContExp}, nil
	case AcceptSpec: return []CompSpec{s.ContExp}, nil
	default: return nil, nil
	}
}
