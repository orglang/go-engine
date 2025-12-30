package typedef

type TypeSpecME struct {
	TypeQN string     `json:"qn"`
	TypeTS TermSpecME `json:"state"`
}

type IdentME struct {
	ID string `json:"id" param:"id"`
}

type TypeRefME struct {
	TypeID string `json:"id" param:"id"`
	TypeRN int64  `json:"rev" query:"rev"`
	Title  string `json:"title"`
}

type TypeSnapME struct {
	TypeID string     `json:"id" param:"id"`
	TypeRN int64      `json:"rev" query:"rev"`
	Title  string     `json:"title"`
	TypeQN string     `json:"qn"`
	TypeTS TermSpecME `json:"state"`
}

type TermSpecME struct {
	K      TermKind    `json:"kind"`
	Link   *LinkSpecME `json:"link,omitempty"`
	Tensor *ProdSpecME `json:"tensor,omitempty"`
	Lolli  *ProdSpecME `json:"lolli,omitempty"`
	Plus   *SumSpecME  `json:"plus,omitempty"`
	With   *SumSpecME  `json:"with,omitempty"`
}

type LinkSpecME struct {
	QN string `json:"qn"`
}

type ProdSpecME struct {
	Value TermSpecME `json:"value"`
	Cont  TermSpecME `json:"cont"`
}

type SumSpecME struct {
	Choices []ChoiceSpecME `json:"choices"`
}

type ChoiceSpecME struct {
	Label string     `json:"label"`
	Cont  TermSpecME `json:"cont"`
}

type TermRefME struct {
	ID string   `json:"id" param:"id"`
	K  TermKind `json:"kind"`
}

type TermKind string

const (
	OneKind    = TermKind("one")
	LinkKind   = TermKind("link")
	TensorKind = TermKind("tensor")
	LolliKind  = TermKind("lolli")
	PlusKind   = TermKind("plus")
	WithKind   = TermKind("with")
)
