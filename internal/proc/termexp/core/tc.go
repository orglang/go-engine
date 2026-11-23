package core

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"fmt"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
)

func marshalTaggedSpec(enc *jsontext.Encoder, kind string, fields map[string]any) error {
	for name, value := range fields {
		if child, ok := value.(ExpSpec); ok && child == nil {
			return fmt.Errorf("%s: nil nested expression", name)
		}
		if branches, ok := value.(map[string]ExpSpec); ok {
			for label, child := range branches {
				if child == nil { return fmt.Errorf("%s[%s]: nil nested expression", name, label) }
			}
		}
	}
	fields["kind"] = kind
	return jsonv2.MarshalEncode(enc, fields)
}

func validateDefinitionSpec(spec ExpSpec) error {
	if spec == nil { return fmt.Errorf("nil process term expression") }
	switch s := spec.(type) {
	case CloseSpec, SendSpec, LinkSpec, FwdSpec, DetachSpec, ReleaseSpec:
		return nil
	case WaitSpec: return validateDefinitionSpec(s.ContExp)
	case RecvSpec: return validateDefinitionSpec(s.ContExp)
	case LabSpec: return validateDefinitionSpec(s.ContExp)
	case CaseSpec:
		for label, child := range s.ContExps {
			if err := validateDefinitionSpec(child); err != nil { return fmt.Errorf("case[%s]: %w", uniqsym.ConvertToString(label), err) }
		}
		return nil
	case CallSpec: return validateDefinitionSpec(s.ContExp)
	case SpawnSpec: return validateDefinitionSpec(s.ContExp)
	case AcqureSpec: return validateDefinitionSpec(s.ContExp)
	case AcceptSpec: return validateDefinitionSpec(s.ContExp)
	default: return ErrExpTypeUnexpected(spec)
	}
}

func (s CloseSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "close", map[string]any{"cont_chnl_ph": string(s.ContChnlPH)})
}
func (s WaitSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "wait", map[string]any{"cont_chnl_ph": string(s.ContChnlPH), "cont_exp": s.ContExp})
}
func (s SendSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "send", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "val_chnl_ph": string(s.ValChnlPH)})
}
func (s RecvSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "recv", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "new_chnl_ph": string(s.NewChnlPH), "cont_exp": s.ContExp})
}
func (s LabSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "lab", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "val_lab_qn": uniqsym.ConvertToString(s.ValLabQN), "cont_exp": s.ContExp})
}
func (s CaseSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	branches := make(map[string]ExpSpec, len(s.ContExps))
	for qn, child := range s.ContExps {
		name := uniqsym.ConvertToString(qn)
		if _, exists := branches[name]; exists { return fmt.Errorf("duplicate case label %q", name) }
		branches[name] = child
	}
	return marshalTaggedSpec(enc, "case", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exps": branches})
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
	chs := make([]string, len(s.ValChnlPHs))
	for i, ch := range s.ValChnlPHs { chs[i] = string(ch) }
	return marshalTaggedSpec(enc, "call", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "new_chnl_ph": string(s.NewChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "val_chnl_phs": chs, "cont_exp": s.ContExp})
}
func (s SpawnSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	chs := make([]string, len(s.NewChnlPHs))
	for i, ch := range s.NewChnlPHs { chs[i] = string(ch) }
	return marshalTaggedSpec(enc, "spawn", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "proc_term_qn": uniqsym.ConvertToString(s.ProcTermQN), "new_chnl_phs": chs, "cont_exp": s.ContExp})
}
func (s AcqureSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "acquire", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp": s.ContExp})
}
func (s AcceptSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "accept", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH), "cont_exp": s.ContExp})
}
func (s DetachSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "detach", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)})
}
func (s ReleaseSpec) MarshalJSONTo(enc *jsontext.Encoder) error {
	return marshalTaggedSpec(enc, "release", map[string]any{"conn_chnl_ph": string(s.ConnChnlPH)})
}
