package typedef

type DefSpecME struct {
	TypeQN string     `json:"type_qn"`
	TypeTS TermSpecME `json:"type_ts"`
}

type IdentME struct {
	DefID string `json:"def_id" param:"id"`
}

type DefRefME struct {
	DefID string `json:"def_id" param:"id"`
	DefRN int64  `json:"def_rn" query:"rn"`
}

type DefSnapME struct {
	DefID  string     `json:"def_id" param:"id"`
	DefRN  int64      `json:"def_rn" query:"rn"`
	Title  string     `json:"title"`
	TypeQN string     `json:"type_qn"`
	TypeTS TermSpecME `json:"type_ts"`
}

type TermSpecME struct {
	K      termKindME  `json:"kind"`
	Link   *LinkSpecME `json:"link,omitempty"`
	Tensor *ProdSpecME `json:"tensor,omitempty"`
	Lolli  *ProdSpecME `json:"lolli,omitempty"`
	Plus   *SumSpecME  `json:"plus,omitempty"`
	With   *SumSpecME  `json:"with,omitempty"`
}

type LinkSpecME struct {
	TypeQN string `json:"type_qn"`
}

type ProdSpecME struct {
	ValTS  TermSpecME `json:"value"`
	ContTS TermSpecME `json:"cont"`
}

type SumSpecME struct {
	Choices []ChoiceSpecME `json:"choices"`
}

type ChoiceSpecME struct {
	Label  string     `json:"label"`
	ContTS TermSpecME `json:"cont"`
}

type TermRefME struct {
	TermID string     `json:"term_id" param:"id"`
	K      termKindME `json:"kind"`
}

type termKindME string

const (
	OneKind    = termKindME("one")
	LinkKind   = termKindME("link")
	TensorKind = termKindME("tensor")
	LolliKind  = termKindME("lolli")
	PlusKind   = termKindME("plus")
	WithKind   = termKindME("with")
)
